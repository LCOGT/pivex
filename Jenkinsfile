#!/usr/bin/env groovy

@Library('lco-shared-libs@0.1.0') _

pipeline {
	agent any
	environment {
		PROJECT_NAME = projName()
	}
	stages {
		stage('Build') {
			steps {
				sh '''
					make
					gzip "${PROJECT_NAME}"
				'''
			}
			post {
			    success { archiveArtifacts artifacts: "${PROJECT_NAME}.gz", fingerprint: true }
			}
		}
		stage('Release') {
			environment {
				GIT_TAG = sh(returnStdout: true, script: "git describe").trim()
				GIT_TAG_MSG = sh(returnStdout: true, script: "git log -1 --pretty=format:%B").trim()
			}
			when { buildingTag() }
			steps {
				githubRelease("${GIT_TAG}", "${GIT_TAG_MSG}", "${PROJECT_NAME}.gz", 'application/gzip')
			}
		}
	}
	post {
		always { postBuildNotify() }
	}
}
