#!/usr/bin/env groovy

@Library('lco-shared-libs@feature/slack') _

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
		success {
			slackSend color: 'good', message: 'Built successfully', channel: '@bkurczynski'
		}
		failure {
			slack()
		}
	}
}
