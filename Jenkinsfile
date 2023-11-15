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
    DOCKERHUB_CREDENTIALS = credentials('dockerhub')
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
        go build -v -o project-privy .
        docker build -t dhutapratama/privy-middleware:latest -f ./Dockerfile .
        """
        // docker build -t dhutapratama/privy-middleware:latest -f ./go-app/Dockerfile .
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
        docker build -t dhutapratama/privy-nodejs-jwt:latest -f ./node-app/Dockerfile .
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
          docker push dhutapratama/privy-middleware:latest
          docker push dhutapratama/privy-nodejs-jwt:latest
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
