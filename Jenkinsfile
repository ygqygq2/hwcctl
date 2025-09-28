pipeline {
  agent any

  parameters {
    text(name: 'URL_LIST', 
      description: '请输入需要刷新的 CDN URL（每行一个）',
      defaultValue: 'https://example.com/index.js',
      trim: true)
  }

  environment {
    HWCCTL_BIN = './bin/hwcctl'
    HWCCTL_CREDENTIALS_ID = 'hwcctl-config'
  }

  stages {
    stage('Validate Parameters') {
      steps {
        script {
          def urls = params.URL_LIST
            .readLines()
            .collect { it.trim() }
            .findAll { it }

          if (urls.isEmpty()) {
            error '未提供任何 URL。'
          }

          env.CDN_URLS = urls.join(',')

          echo "准备刷新 ${urls.size()} 个 URL：\n - ${urls.join('\n - ')}"
        }
      }
    }

    stage('Refresh CDN') {
      steps {
        withCredentials([
          file(credentialsId: env.HWCCTL_CREDENTIALS_ID, variable: 'HWCCTL_CONFIG_FILE')
        ]) {
          sh '''
            set -euo pipefail
            echo "使用配置文件刷新 CDN，URL 列表：$CDN_URLS"
            $HWCCTL_BIN --config "$HWCCTL_CONFIG_FILE" cdn refresh --urls "$CDN_URLS"
          '''
        }
      }
    }
  }
}
