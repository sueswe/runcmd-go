pipeline {
    agent any
    environment {
        J_URL = "${env.JENKINS_URL}"
        J_JOB = "${env.JOB_NAME}"
        J_BUILD = "${BUILD_NUMBER}"
        J_BUILD_URL = "${BUILD_URL}"
        J_JOB_URL = "${JOB_URL}"
        BRANCHNAME = "${BRANCH_NAME}"
    }

    stages {
        stage('git clone') {
            when {
                branch "master"
            }
            steps {
                script {
                    try {
                        sh ' rm -rf ~/compile/runcmd-go '
                        sh ' cd ~/compile/ && git clone git@lvgom01.sozvers.at:repositories/runcmd-go.git '
                    } catch(err) {
                        throw err
                    }
                }
            }
        }
        stage('compile') {
            steps {
                script {
                    try {
                        sh ' cd ~/compile/runcmd-go && .deploy/compile.sh '
                    } catch(err) {
                        throw err
                    }
                }
            }
        }
        stage('deploy') {
            steps {
                script {
                    try {
                        sh ' cd ~/compile/runcmd-go && .deploy/deploy.sh '
                    } catch(err) {
                        throw err
                    }
                }
            }
        }
    }
}
