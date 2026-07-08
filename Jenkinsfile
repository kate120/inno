pipeline {
    agent any

        stages {

            stage('Setup Python') {
                steps {
                    dir('module-08') {
                        sh 'pip install -r requirements.txt'
                    }
                }
            }

            stage('Test') {
                steps {
                    dir('module-08') {
                        sh 'pytest tests/ -v'
                    }
                }
            }


        }
}