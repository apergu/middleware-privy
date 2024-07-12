properties(
[pipelineTriggers([pollSCM('* * * * *')])]
)

def FAILED_STAGE

pipeline {
  agent any

  //environment
  environment {
    // Repository
    def GIT_CREDENTIAL = "git.dev1.my.id"
    def GIT_HASH = sh(returnStdout: true, script: 'git log -1 --pretty=format:"%h"').trim()
    DOCKERHUB_CREDENTIALS = credentials('dockerhub-apergu')
  }

  stages {
    stage("BUILD MIDDLEWARE") {
      steps {
        script {
          FAILED_STAGE=env.STAGE_NAME
          echo "BUILD MIDDLEWARE"
        }
        sh label: 'Build Middleware Script', script:
        """
        docker build -t apergudev/privy-middleware:latest -f ./Dockerfile .
        docker build -t apergudev/privy-middleware:dev -f ./Dockerfile .
        docker build -t apergudev/privy-middleware:staging -f ./Dockerfile .
        """
        // go build -v -o privy .
        // docker build -t apergudev/privy-middleware:latest -f ./go-app/Dockerfile .
      }
    }    
    
    stage("BUILD NODEJS") {
      steps {
        script {
          FAILED_STAGE=env.STAGE_NAME
          echo "BUILD NODEJS"
        }
        sh label: 'Build NodeJS-jwt Script', script:
        """
        docker build -t apergudev/privy-nodejs-jwt:latest -f ./node-app/Dockerfile .
        docker build -t apergudev/privy-nodejs-jwt:dev -f ./node-app/Dockerfile .
        docker build -t apergudev/privy-nodejs-jwt:staging -f ./node-app/Dockerfile .
        """
      }
    }

    stage("RELEASE") {
      steps {
        script {
          FAILED_STAGE=env.STAGE_NAME
          echo "RELEASE"
        }

        sh label: 'STEP RELEASE', script:
        """
          echo $DOCKERHUB_CREDENTIALS_PSW | docker login -u $DOCKERHUB_CREDENTIALS_USR --password-stdin
          docker push apergudev/privy-middleware:latest
          docker push apergudev/privy-middleware:dev
          docker push apergudev/privy-middleware:staging
          docker push apergudev/privy-nodejs-jwt:latest
          docker push apergudev/privy-nodejs-jwt:dev
          docker push apergudev/privy-nodejs-jwt:staging
        """
      }
    }
  }
  post {
    always {
      sh 'docker logout'
    }
  }
}
