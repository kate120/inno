pipeline {
    agent any

    stages {

        stage('Setup Python') {
            steps {
                dir('module-08') {
                    sh '''
                        python3 -m venv venv
                        venv/bin/pip install --no-cache-dir -r requirements.txt
                    '''
                }
            }
        }

        stage('Test') {
            steps {
                dir('module-08') {
                    sh 'venv/bin/pytest tests/ -v'
                }
            }
        }

    }
}