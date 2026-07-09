pipeline {
    agent any

    stages {

        stage('Check branch') {
            steps {
                script {
                    def branch = env.BRANCH_NAME
                    echo "this branch is ${branch}"
                }
            }
        }

        stage('chec files') {
            steps {
                echo "files in the branch"
                sh "ls"
            }
        }

        stage('PR') {
            when {
                expression { env.CHANGE_ID != null }  // Только для PR
            }
            steps {
                echo "PR has been checked"
            }
        }

    }


}