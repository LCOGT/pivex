#!/usr/bin/env groovy

@Library('lco-shared-libs@0.1.0') _

pipeline {
	agent any
	stages {
		stage('Build') {
			steps {
			// TODO: gzip the binary
				sh 'make'
			}
			post {
			    success { archiveArtifacts artifacts: "pivex.gz", fingerprint: true }
			}
		}
		stage('Release') {
			steps {
				githubRelease('tag', 'description', 'filename', 'application/gzip')
			}
		}
	}
	post {
		always { postBuildNotify() }
	}
}
