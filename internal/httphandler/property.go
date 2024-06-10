package httphandler

import (
	"middleware/pkg/credential"
	"middleware/pkg/erpprivy"
	"middleware/pkg/privy"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/rteja-library3/rcache"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/remailer"
	"gitlab.com/rteja-library3/rpassword"
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
	DefaultCredential   credential.Credential
	DefaultPrivy        privy.Privy
	DefaultERPPrivy     erpprivy.ErpPrivy
}
