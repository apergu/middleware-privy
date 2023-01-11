const express = require('express')
const app = express()
const port = 3000

var rs = require('jsrsasign');

app.get('/', (req, res) => {
    var jwtHeader = {
        alg: 'PS256',
        typ: 'JWT',
        kid: 'T17p9_H2k1WjJB3h3vQk0qFJlbFERxUZRx2w71TDw7k' 
    };
    
    let stringifiedJwtHeader = JSON.stringify(jwtHeader);

    let jwtPayload = {
        iss: 'f8b4fb43614c756a168b1d1367607569ae547e8ffb6001b0ec7526641627c112', 
        scope: ['restlets','rest_webservices'], 
        iat: (new Date() / 1000),               
        exp: (new Date() / 1000) + 3600,        
        aud: 'https://tstdrv2245019.suitetalk.api.netsuite.com/services/rest/auth/oauth2/v1/token'
    };

    var stringifiedJwtPayload = JSON.stringify(jwtPayload);
    
    let secret = '-----BEGIN PRIVATE KEY-----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-----END PRIVATE KEY-----';

    var JWS = rs.jws.JWS;

    let signedJWT = JWS.sign('PS256',stringifiedJwtHeader,stringifiedJwtPayload,secret);

    res.setHeader('Content-Type', 'application/json');
    res.end(JSON.stringify({ signedJWT }));
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})