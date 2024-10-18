const express = require('express')
const winston = require('winston')
const { format } = require('date-fns')
const app = express()
const port = 3000

// Function to create a log filename based on the current date
const getLogFilename = () => {
    const now = new Date();
    return `application-${format(now, 'yyyy-MM-dd-HH:mm')}.log`; // Create the filename
};


// Set up the logger
const logger = winston.createLogger({
    level: 'info',
    format: winston.format.combine(
        winston.format.timestamp(),
        winston.format.json()
    ),
    transports: [
        new winston.transports.File({ filename: getLogFilename() }) // Log file will be created in the same folder
    ],
});



var rs = require('jsrsasign');

app.get('/', (req, res) => {
    logger.info('Received a request to the root endpoint on time: ' + new Date().toISOString());

    var jwtHeader = {
        alg: 'PS256',
        typ: 'JWT',
        kid: 'r49s5LIxZvIAMvKbep8p-kNphXgdNljdI_0rNhK5axk' 
    };
    
    let stringifiedJwtHeader = JSON.stringify(jwtHeader);

    let jwtPayload = {
        iss: 'd0e9d94e53648143208032f312cbe8c47cfa1a2084228114bcb4e2ce397f8704', 
        scope: ['restlets','rest_webservices'], 
        iat: (new Date() / 1000),               
        exp: (new Date() / 1000) + 3600,        
        aud: 'https://8113915-sb1.suitetalk.api.netsuite.com/services/rest/auth/oauth2/v1/token'
    };
    
    var stringifiedJwtPayload = JSON.stringify(jwtPayload);

    logger.info('Creating JWT with header and payload', { header: jwtHeader, payload: jwtPayload });

    
    // let secret = '-----BEGIN PRIVATE KEY-----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-----END PRIVATE KEY-----';

    let secret = '-----BEGIN PRIVATE KEY-----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-----END PRIVATE KEY-----';
    var JWS = rs.jws.JWS;

    let signedJWT = JWS.sign('PS256',stringifiedJwtHeader,stringifiedJwtPayload,secret);

    res.end(JSON.stringify({ signedJWT }));

    logger.info('============================= END ==================== ');
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})