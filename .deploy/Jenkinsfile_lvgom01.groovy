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
        
        stage('LVGOM01:') {
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
                script {
                    try {
                        sh ' cd ~/compile/runcmd-go && .deploy/compile_deploy.sh '
                    } catch(err) {
                        throw err
                    }
                }
            }
        }
    }
}

