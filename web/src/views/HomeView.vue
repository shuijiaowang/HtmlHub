<template>
  <div class="home-page">
    <main class="content">
      <!-- 平台头部介绍 -->
      <section class="card hero-card">
        <h1>HtmlHub</h1>
        <p class="subtitle">一键上传 HTML，分配子域名，轻松托管静态页面</p>
        <div class="hero-tags">
          <span>零部署</span>
          <span>多设备存储同步</span>
          <span>私密/公开访问</span>
          <span>一键分享</span>
        </div>
      </section>

      <!-- 核心功能 -->
      <section class="card">
        <h2>核心功能</h2>
        <div class="feature-grid">
          <div class="feature-item">
            <h4>简单上传</h4>
            <p>支持上传 HTML 文件 / 直接粘贴代码，无需复杂配置</p>
          </div>
          <div class="feature-item">
            <h4>自定义子域名</h4>
            <p>自定义子域名，自动分配唯一域名，专属访问链接，方便分享</p>
          </div>
          <div class="feature-item">
            <h4>本地存储同步</h4>
            <p>内置 localStorage 云端同步脚本，支持数据导出/导入/上传/加载</p>
          </div>
          <div class="feature-item">
            <h4>权限控制</h4>
            <p>页面支持私密/公开切换，未公开页面仅自己可访问，保护隐私</p>
          </div>
          <div class="feature-item">
            <h4>数据统计</h4>
            <p>统计页面访问次数，清晰掌握页面使用情况</p>
          </div>
          <div class="feature-item">
            <h4>完整管理</h4>
            <p>管理所有页面，支持编辑、更新、删除、修改描述等操作</p>
          </div>
        </div>
      </section>

      <!-- 快速上手教程 -->
      <section class="card">
        <h2>快速上手</h2>
        <div class="guide-steps">
          <div class="step">
            <span class="step-num">1</span>
            <div class="step-content">
              <h4>上传 HTML</h4>
              <p>进入「上传 HTML」页面，填写子域名（可选），上传文件或粘贴完整 HTML 代码</p>
            </div>
          </div>
          <div class="step">
            <span class="step-num">2</span>
            <div class="step-content">
              <h4>获取访问链接</h4>
              <p>上传成功后，系统自动分配子域名，生成专属页面访问链接</p>
            </div>
          </div>
          <div class="step">
            <span class="step-num">3</span>
            <div class="step-content">
              <h4>管理与分享</h4>
              <p>在「个人 HTML 管理」中查看、编辑页面，设置公开/私密，分享给他人</p>
            </div>
          </div>
        </div>
      </section>

      <!-- AI 创作提示 -->
      <section class="card">
        <h2>AI 快速生成页面</h2>
        <p class="tips-title">复制下方提示词，直接发给 AI，一键生成可托管的完整 HTML 页面：</p>
        <div class="code-block">
          请用html+css+js+local存储生成一个完整的html页面，需求是：...
        </div>
        <p class="tips-desc">生成完成后，直接复制代码上传到本平台即可使用</p>
      </section>

      <!-- 联系与支持 -->
      <section class="card">
        <h2>联系站长</h2>
        <p class="platform-desc">
          本网站主要为非程序员打造，不用自己部署服务器、无需搭建数据库，解决纯静态网页无法数据同步、无法分享页面的问题，平台目前尚处于开发阶段，请注意备份自己的数据。<br>
          ✅ 如需纯净子域名、修改使用上限，可联系站长；<br>
          ✅ 如遇到网站bug、功能建议，或是 VibeCoding 创作过程中遇到问题，均可联系站长帮忙解答。
        </p>
        <div class="contact-list">
          <div class="contact-item">
            <span>QQ 群：</span>
            <span>1102036440</span>
          </div>
          <div class="contact-item">
            <span>微信号：</span>
            <span>shuijiaowang1</span>
          </div>
          <div class="contact-item">
            <span>独立部署：</span>
            <span>支持个人网站/服务器/数据库定制部署服务</span>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
.home-page {
  min-height: 100vh;
  background:
    radial-gradient(900px 420px at 20% -10%, rgb(var(--hh-brand-rgb) / 0.12), transparent 60%),
    radial-gradient(900px 420px at 100% 0%, rgb(var(--hh-brand-2-rgb) / 0.10), transparent 55%),
    var(--color-background);
}

.content {
  padding: 36px 24px;
  max-width: 1200px;
  margin: 0 auto;
}

/* 通用卡片样式 */
.card {
  background: var(--hh-surface-solid);
  border: 1px solid var(--hh-border);
  border-radius: 16px;
  padding: 26px;
  margin-bottom: 24px;
  box-shadow: var(--hh-shadow-md);
}

/* 头部英雄区 */
.hero-card {
  text-align: center;
  padding: 44px 24px;
  background:
    radial-gradient(650px 280px at 20% 10%, rgb(var(--hh-brand-rgb) / 0.16), transparent 55%),
    radial-gradient(650px 280px at 90% 0%, rgb(var(--hh-brand-2-rgb) / 0.14), transparent 55%),
    linear-gradient(135deg, #f7fbff 0%, #f3f8ff 100%);
}

.hero-card h1 {
  font-size: clamp(28px, 4.2vw, 40px);
  font-weight: 850;
  color: var(--color-heading);
  margin-bottom: 10px;
  letter-spacing: -0.3px;
}

.subtitle {
  font-size: 15px;
  color: var(--hh-text-2);
  margin-bottom: 20px;
}

.hero-tags {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 10px;
}

.hero-tags span {
  padding: 7px 14px;
  background: rgb(var(--hh-brand-rgb) / 0.10);
  color: color-mix(in srgb, var(--hh-brand) 88%, #000 0%);
  border-radius: 20px;
  font-size: 13px;
  border: 1px solid rgb(var(--hh-brand-rgb) / 0.16);
}

/* 标题统一样式 */
.card h2 {
  font-size: 22px;
  font-weight: 800;
  color: var(--color-heading);
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgb(var(--hh-brand-rgb) / 0.18);
  letter-spacing: -0.2px;
}

/* 功能网格 */
.feature-grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: 20px;
}

.feature-item {
  grid-column: span 4;
  padding: 18px;
  border-radius: 14px;
  background:
    radial-gradient(420px 140px at 20% 0%, rgb(var(--hh-brand-rgb) / 0.12), transparent 60%),
    color-mix(in srgb, var(--hh-surface-solid) 94%, var(--hh-brand) 0%);
  border: 1px solid rgba(2, 6, 23, 0.06);
  box-shadow: 0 10px 24px rgba(2, 6, 23, 0.06);
}

.feature-item h4 {
  font-size: 16px;
  font-weight: 800;
  color: var(--color-heading);
  margin-bottom: 8px;
}

.feature-item p {
  font-size: 14px;
  color: var(--hh-text-2);
  line-height: 1.55;
  margin: 0;
}

/* 步骤引导 */
.guide-steps {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.step {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 14px 14px;
  border-radius: 14px;
  border: 1px solid rgba(2, 6, 23, 0.06);
  background: color-mix(in srgb, var(--hh-surface-solid) 96%, #000 0%);
}

.step-num {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: color-mix(in srgb, var(--hh-brand) 86%, #000 0%);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  flex-shrink: 0;
}

.step-content h4 {
  font-size: 16px;
  font-weight: 800;
  color: var(--color-heading);
  margin: 0 0 4px 0;
}

.step-content p {
  font-size: 14px;
  color: var(--hh-text-2);
  margin: 0;
  line-height: 1.5;
}

/* AI 提示代码块 */
.tips-title {
  font-size: 15px;
  color: var(--hh-text-2);
  margin-bottom: 12px;
}

.code-block {
  padding: 16px;
  background: rgba(2, 6, 23, 0.04);
  border-radius: 12px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  color: var(--hh-text);
  border: 1px solid rgba(2, 6, 23, 0.06);
  margin-bottom: 8px;
  line-height: 1.6;
  overflow: auto;
}

.tips-desc {
  font-size: 14px;
  color: var(--hh-text-3);
  margin: 0;
}

/* 新增：平台说明文字样式 */
.platform-desc {
  font-size: 15px;
  color: var(--hh-text-2);
  line-height: 1.6;
  margin-bottom: 20px;
  padding: 14px;
  background: rgb(var(--hh-brand-rgb) / 0.06);
  border-radius: 12px;
  border: 1px solid rgb(var(--hh-brand-rgb) / 0.14);
}

/* 联系方式 */
.contact-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.contact-item {
  display: flex;
  gap: 8px;
  font-size: 15px;
}

.contact-item span:first-child {
  color: var(--color-heading);
  font-weight: 700;
  min-width: 80px;
}

.contact-item span:last-child {
  color: var(--hh-text-2);
}

@media (max-width: 640px) {
  .content {
    padding: 18px 12px;
  }

  .card {
    padding: 16px;
    border-radius: 14px;
  }

  .feature-grid {
    gap: 14px;
  }

  .feature-item {
    grid-column: span 12;
  }

  .step {
    padding: 12px;
    gap: 12px;
  }
}

@media (max-width: 1024px) and (min-width: 641px) {
  .feature-item {
    grid-column: span 6;
  }
}
</style>