package httphandler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/mohamadikbal/project-privy/pkg/credential"
	"gitlab.com/rteja-library3/rcache"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/remailer"
	"gitlab.com/rteja-library3/rpassword"
	"gitlab.com/rteja-library3/rstorager"
	"gitlab.com/rteja-library3/rtoken"
)

type HTTPHandlerProperty struct {
	DBPool              *pgxpool.Pool
	DefaultDecoder      rdecoder.Decoder
	DefaultEmailer      remailer.Remail
	DefaultPwdEncryptor rpassword.Encryptor
	DefaultCache        rcache.Cache
	DefaultToken        rtoken.Token
	DefaultRefreshToken rtoken.Token
	DefaultStorage      rstorager.Storage
	DefaultCredential   credential.Credential
}
