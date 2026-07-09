pipeline {
    agent any

    environment {
        IMAGE_NAME = "kiwwa/my-app-jenkins"
        IMAGE_TAG  = "${env.BUILD_NUMBER}"
    }

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

        stage('Build Image') {
            steps {
                dir('module-08') {
                    sh "docker build -t ${IMAGE_NAME}:${IMAGE_TAG} ."
                }
            }
        }

        stage('Push to Docker Hub') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-creds',
                    usernameVariable: 'DOCKER_USERNAME',
                    passwordVariable: 'DOCKER_PASSWORD'
                )]) {
                    sh '''
                        echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
                        docker push ${IMAGE_NAME}:${IMAGE_TAG}
                        docker logout
                    '''
                }
            }
        }

        stage('PR') {
            when {
                expression { env.CHANGE_ID != null }
            }
            steps {
                echo "PR has been checked"
            }
        }

    }
}