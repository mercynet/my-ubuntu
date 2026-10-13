"""Preparation regressions use synthetic files and stubbed subprocesses only."""
import copy
import importlib.util
import json
from pathlib import Path
import tempfile
import threading
import time
import unittest
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[2]
loader = importlib.util.spec_from_file_location('prepare', ROOT/'internal/integrations/prepare-compose.py')
setup = importlib.util.module_from_spec(loader)
loader.loader.exec_module(setup)


class PreparationTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.project = Path(self.temp.name)/'project'
        self.project.mkdir()
        self.workspace = Path(self.temp.name)/'workspace'
        self.spec = {'compose':'docker-compose.yml','application':'app','worker':'worker','initializer':'redis-cluster-init','template':'docker/bin/tenant-create.sh','databases':[{'service':'mariadb','schemas':['tenant1']},{'service':'mariadb2','schemas':['tenant2']}],'redis':[{'service':'redis-node-'+str(i),'port':7000+i} for i in range(1,4)],'tenants':[{'domain':'og-tenant1.test','service':'mariadb','schema':'tenant1'}]}
        (self.project/'docker-compose.yml').write_text('original')
        template = self.project/self.spec['template']
        template.parent.mkdir(parents=True)
        template.write_text('cat > "${CONFIG_DIR}/config.php" <<EOF\n<?php\nconst DB_SERVER = \'${SELECTED_CONTAINER}\';\nconst DB_PASS = \'${DB_PASS}\';\nconst DB_DATABASE = \'${NEW_SCHEMA}\';\nconst REDIS_HOSTS = "redis:6379";\nEOF\n')
        services = {name: {'image': 'fixture', 'networks': {'local': {}}, 'restart': 'always'}
                    for name in ['app','worker','nginx','redis-cluster-init','mariadb','mariadb2','redis-node-1','redis-node-2','redis-node-3']}
        services['app']['build'] = {'context': str(self.project)}
        services['nginx']['ports'] = [{'target': 80}, {'target': 443}]
        services['mariadb']['volumes'] = [{'type':'bind','source':str(self.project/'init'), 'target':'/docker-entrypoint-initdb.d'}]
        self.cfg = {'services':services,'networks':{'local':{'external':True}},'volumes':{'db':{'external':True}}}
        self.details = {'cfg':self.cfg, 'settings':{}, 'namespace':'mu-fixture', 'ports':{'http':8089,'https':8449},'steps':['configuration','build','up','database','redis','composer','verify'], 'manual':[], 'address':'https://localhost:8449'}

    def fake_run(self, argv, *args, **kwargs):
        if argv[:2] == ['openssl','req']:
            Path(argv[argv.index('-keyout')+1]).write_text('synthetic key')
            Path(argv[argv.index('-out')+1]).write_text('synthetic certificate')
        return ''

    def prepare(self):
        with patch.object(setup,'run',side_effect=self.fake_run):
            return setup.prepare(self.project,self.spec,self.workspace,self.details)

    def test_isolation_preserves_sources_and_existing_tenant(self):
        tenant = self.project/'_configs/og-tenant1.test:8089/config.php'
        tenant.parent.mkdir(parents=True)
        tenant.write_text('existing configuration')
        before = copy.deepcopy(self.cfg)
        cfg = self.prepare()
        self.assertEqual(self.cfg,before)
        self.assertEqual((self.project/'docker-compose.yml').read_text(),'original')
        self.assertEqual(tenant.read_text(),'existing configuration')
        self.assertEqual(cfg['services']['worker']['image'],cfg['services']['app']['image'])
        for name,service in cfg['services'].items():
            self.assertEqual(service['container_name'],'mu-fixture-'+name)
            self.assertEqual(service['labels'][setup.LABEL],'mu-fixture')
            self.assertEqual(service['restart'],'no')
            for port in service.get('ports',[]):self.assertEqual(port['host_ip'],'127.0.0.1')
        for kind in ('networks','volumes'):
            for name,value in cfg[kind].items():
                self.assertNotIn('external',value)
                self.assertEqual(value['name'],'mu-fixture_'+name)
        script=(self.workspace/'init-cluster.sh').read_text()
        self.assertNotIn('reset',script.lower())
        self.assertIn('cluster_slots_assigned:0',script)
        self.assertIn('--cluster-replicas 0',script)
        sql=(self.workspace/'mariadb/00-init.sql').read_text()
        self.assertIn('CREATE DATABASE IF NOT EXISTS',sql)
        self.assertNotIn('DROP',sql)
        settings=json.loads((self.workspace/'settings.json').read_text())
        self.details['settings']=settings
        self.prepare()
        self.assertEqual(json.loads((self.workspace/'settings.json').read_text())['password'],settings['password'])

    def test_symlink_and_escape_are_refused(self):
        outside=Path(self.temp.name)/'outside'
        outside.write_text('sentinel')
        with self.assertRaises(ValueError):setup.inside(self.project,'../outside')
        self.workspace.symlink_to(self.project,target_is_directory=True)
        with self.assertRaises(ValueError):self.prepare()
        self.assertEqual(outside.read_text(),'sentinel')

    def test_resume_refuses_changed_project(self):
        control=str(Path(self.temp.name)/'control')
        flag=Path(control+'.pause');flag.touch()
        events=[]
        errors=[]
        with patch.object(setup,'emit',side_effect=lambda **values:events.append(values)):
            def wait():
                try:setup.boundary(control,self.project,self.spec)
                except Exception as error:errors.append(error)
            worker=threading.Thread(target=wait);worker.start()
            deadline=time.monotonic()+2
            while not events and time.monotonic()<deadline:time.sleep(.01)
            self.assertTrue(events)
            (self.project/'docker-compose.yml').write_text('changed')
            flag.unlink();worker.join(2)
        self.assertFalse(worker.is_alive())
        self.assertIsInstance(errors[0],ValueError)

    def test_cancel_blocks_next_operation(self):
        control=str(Path(self.temp.name)/'control')
        Path(control+'.cancel').touch()
        with self.assertRaises(InterruptedError):setup.boundary(control,self.project,self.spec)

    def test_lock_blocks_second_preparation_and_releases(self):
        path=Path(self.temp.name)/'setup.lock'
        with setup.acquire_lock(path):
            with self.assertRaises(ValueError):setup.acquire_lock(path)
        with setup.acquire_lock(path):pass

    def test_missing_or_unhealthy_services_never_report_success(self):
        prefix=['fixture']
        with patch.object(setup,'run',return_value='[]'):
            with self.assertRaises(ValueError):setup.verify_services(prefix,self.cfg,self.spec)
        values=[{'Service':name,'State':'running','Health':'healthy'} for name in self.cfg['services']]
        for v in values:
            if v['Service']=='redis-cluster-init':v.update(State='exited',ExitCode=0)
        with patch.object(setup,'run',return_value=json.dumps(values)):
            setup.verify_services(prefix,self.cfg,self.spec)
        values[0]['Health']='unhealthy'
        with patch.object(setup,'run',return_value=json.dumps(values)):
            with self.assertRaises(ValueError):setup.verify_services(prefix,self.cfg,self.spec)


if __name__ == '__main__':
    unittest.main()
