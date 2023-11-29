const express = require('express')
const app = express()
const port = 3000

var rs = require('jsrsasign');

app.get('/', (req, res) => {
    var jwtHeader = {
        alg: 'PS256',
        typ: 'JWT',
        kid: 'wS8R74-PHvXF-WMXNiNKZ3WOif_qwSYtbpvXav_9eM8'
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