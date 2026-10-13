# Pacote OfficeGest para testes

Versão 0.2.0, formato 2; exige Meu Ubuntu 1.0.0-rc.4. Exemplo público e sanitizado; não inclui o checkout OfficeGest.
O ZIP contém apenas manifest.json: IDs, serviços, schemas, domínios e guias.
Não inclui executáveis, credenciais, dumps ou fontes OfficeGest.

Importe o ZIP em Pacotes de integração; escolha um checkout de desenvolvimento
com .env sanitizado, composer.lock e o modelo de tenant fornecidos pela equipa.
Reveja ferramentas necessárias e aplique o plano. Depois reveja a preparação
do projeto e confirme. Os comandos da automação pertencem ao Meu Ubuntu.

A preparação cria uma configuração Compose pessoal com nomes e volumes
isolados, portas em 127.0.0.1, credenciais locais, schemas, Redis Cluster e
Composer. O Compose original, tenants existentes e dados são preservados.
O inicializador Redis recusa estado parcial sem executar reset. Containers
não arrancam automaticamente no reboot. Pode pausar/retomar entre etapas;
configurações alteradas durante a pausa exigem nova revisão.

Dumps, /etc/hosts, confiança do certificado e autenticação privada Composer
continuam manuais. O comando dock e as instruções aparecem no resultado.
Se o computador ou aplicação encerrar, reveja novamente para conferir o estado.

Testes sintéticos e simulação não comprovam arranque completo OfficeGest.
Validar num desktop descartável com checkout sanitizado e dumps da equipa
antes de considerar o ambiente pronto. Não iniciar esta automação sobre o
ambiente de trabalho existente para testar o addon.

Gerar o pacote e o ZIP de teste depois do build:

```sh
python3 build/scripts/integration-test-bundle.py --manifest examples/integrations/officegest/manifest.json --version 1.0.0-rc.4
```
