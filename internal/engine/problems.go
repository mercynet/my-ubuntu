package engine

import (
	"strings"

	"meuubuntu/internal/store"
)

// Only recognise concrete diagnostics. An exit code alone cannot identify a cause.
func executionProblem(message, output string) *store.Problem {
	s := strings.ToLower(message + "\n" + output)
	has := func(parts ...string) bool {
		for _, p := range parts {
			if strings.Contains(s, p) {
				return true
			}
		}
		return false
	}
	switch {
	case has("outra sessão administrativa"):
		return &store.Problem{Code: "busy", Summary: "Outra sessão administrativa do Meu Ubuntu está em curso.", Advice: "Aguarde que essa sessão termine e volte a rever o plano."}
	case has("could not get lock", "unable to acquire the dpkg", "is another process using it"):
		return &store.Problem{Code: "busy", Summary: "Outra operação está a usar o gestor de pacotes.", Advice: "Aguarde o fim do Atualizador ou da outra instalação e tente novamente."}
	case has("no space left on device", "disk quota exceeded", "sem espaço no dispositivo"):
		return &store.Problem{Code: "space", Summary: "Não há espaço disponível para concluir a operação.", Advice: "Verifique o espaço livre e as quotas do disco; liberte espaço antes de tentar novamente."}
	case has("temporary failure resolving", "could not resolve", "couldn't resolve", "network is unreachable", "connection timed out", "context deadline exceeded"):
		return &store.Problem{Code: "network", Summary: "Não foi possível contactar o serviço ou resolver o endereço.", Advice: "Verifique a ligação, o DNS e eventuais configurações de VPN/proxy. Depois tente novamente."}
	case has("conflicting values set for option", "the list of sources could not be read", "malformed entry"):
		return &store.Problem{Code: "sources", Summary: "As fontes APT têm uma configuração inválida.", Advice: "Consulte no registo o ficheiro indicado. Corrija a fonte e regenere o plano antes de tentar novamente."}
	case has("no_pubkey", "badsig", "expkeysig", "is not signed", "fingerprint não", "impressão digital não", "fingerprint invál", "impressão digital invál", "sem impressão digital verificada", "chave sem impressões digitais", "chave do repositório não autorizada"):
		return &store.Problem{Code: "signature", Summary: "Não foi possível verificar a identidade do repositório.", Advice: "Confirme a chave e a fonte junto do fornecedor; mantenha a verificação de assinaturas ativa."}
	case has("dpkg was interrupted"):
		return &store.Problem{Code: "dpkg", Summary: "Há uma configuração de pacotes interrompida.", Advice: "Conclua a configuração pendente com sudo dpkg --configure -a no terminal e depois regenere o plano."}
	case has("packages need to be removed but remove is disabled"):
		return &store.Problem{Code: "removal-blocked", Summary: "Esta instalação foi bloqueada para evitar remover outros pacotes.", Advice: "Reveja o conflito indicado no registo. Ajuste as escolhas ou a fonte do pacote e volte a rever o plano."}
	case has("unable to locate package", "has no installation candidate"):
		return &store.Problem{Code: "package", Summary: "O pacote não foi encontrado nos índices disponíveis.", Advice: "Verifique as fontes e atualize os índices pelo Atualizador. Depois regenere o plano para verificar a disponibilidade."}
	case has("autenticação recusada", "autenticação cancelada", "not authorized", "authorization failed"):
		return &store.Problem{Code: "authorization", Summary: "A autorização administrativa não foi concluída.", Advice: "Tente novamente e autorize a operação na janela do sistema quando esta aparecer."}
	case has("script mudou desde a revisão", "plano mudou"):
		return &store.Problem{Code: "review", Summary: "O conteúdo mudou desde a revisão.", Advice: "Gere e reveja um novo plano antes de executar."}
	default:
		return &store.Problem{Code: "unknown", Summary: "Não foi possível concluir ou confirmar este passo.", Advice: "Abra os detalhes e consulte o erro e o registo antes de decidir como continuar."}
	}
}
