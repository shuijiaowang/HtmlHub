
<template>
  <div class="home-page">
    <main class="content">
      <section class="card">
        <h1>HtmlHub 主页</h1>
        <p class="subtitle">一个简单的 HTML 上传与管理平台。</p>
        <h3>使用教程</h3>
        <p>1. 进入「上传 HTML」页面，可填写页面前缀，上传 html 文件或直接粘贴 html 文本。</p>
        <p>2. 点击「上传」后，系统会自动生成文件名并写入记录。</p>
        <p>3. 进入「个人 HTML 管理」可查看你上传过的记录与分享链接。</p>
        <h3>如何问 AI</h3>
        <p class="tips">
          请用html+css+js+local存储生成一个完整的html页面，需求是：...
        </p>
        <h3>联系方式</h3>
        <p>QQ群：待补充</p>
        <p>微信：待补充</p>
      </section>
    </main>
  </div>
</template>

<style scoped>
.home-page {
  min-height: 100vh;
  background: var(--color-background);
}

.content {
  padding: 32px 24px;
  max-width: 980px;
  margin: 0 auto;
}

.subtitle {
  color: #666;
  margin-bottom: 20px;
}

.card {
  background: #fff;
  border: 1px solid var(--color-border);
  border-radius: 10px;
  padding: 18px;
  margin-bottom: 18px;
}

.upload-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.tips {
  margin: 0 0 10px;
  color: #666;
  font-size: 14px;
}
</style>