#!/usr/bin/env groovy

@Library('lco-shared-libs@0.0.4') _

pipeline {
	agent any
	options {
		timeout(time: 10, unit: 'MINUTES')
	}
	environment {
		GOPATH = "${WORKSPACE}"
		PROJ_NAME = projName()
	}
	stages {
		stage('Build') {
			steps {
				sh '''
					mkdir -p src/${PROJ_NAME}
					mv $(ls | grep -v 'src') src/${PROJ_NAME}
				'''

				dir("src/${PROJ_NAME}") {
					sh 'go get && go install'
				}
			}
		}
	}
	post {
		always {
			slack()
		}
	}
}
