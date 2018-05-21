#!/usr/bin/env groovy

@Library('lco-shared-libs@0.0.4') _

pipeline {
	agent any
	options {
		timeout(time: 10, unit: 'MINUTES')
	}
	stages {
		stage('Build') {
			steps {
				sh 'go build'
			}
		}
	}
	post {
		slack()
	}
}
