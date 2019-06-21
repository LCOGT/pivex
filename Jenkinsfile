#!/usr/bin/env groovy

@Library('lco-shared-libs@0.0.5') _

pipeline {
	agent any
	stages {
		stage('Build') {
			steps {
				sh 'make'
			}
		}
		stage('Release') {
			steps {
				sh ':'
			}
		}
	}
	post {
		always {
			slack()
		}
	}
}
