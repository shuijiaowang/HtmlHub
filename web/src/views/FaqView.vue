<template>
  <div class="help-page">
    <main class="content">
      <section class="card hero-card">
        <p class="eyebrow">帮助中心</p>
        <h1>HtmlHub 怎么用？</h1>
        <p class="subtitle">这里用普通话解释常见功能和提示，遇到看不懂的地方可以先来这里找答案。</p>
        <div class="quick-links" aria-label="快速查看">
          <a href="#start">新手上手</a>
          <a href="#visibility">公开 / 私密</a>
          <a href="#publish">发布模式</a>
          <a href="#problems">常见问题</a>
        </div>
      </section>

      <section id="start" class="card">
        <h2>第一次使用，按这几步来</h2>
        <div class="steps">
          <article class="step">
            <span class="step-num">1</span>
            <div>
              <h3>先准备一个 HTML 页面</h3>
              <p>可以让 AI 帮你生成，也可以上传自己写好的 HTML 文件。建议上传完整页面，不要只上传一小段代码。</p>
            </div>
          </article>
          <article class="step">
            <span class="step-num">2</span>
            <div>
              <h3>到「上传」提交页面</h3>
              <p>可以填写一个好记的子域名前缀，例如 <strong>my-game</strong>。如果不填，系统会自动生成。</p>
            </div>
          </article>
          <article class="step">
            <span class="step-num">3</span>
            <div>
              <h3>到「管理」打开和修改</h3>
              <p>上传后，在「管理」里可以打开页面、编辑简介、更新 HTML、删除页面、切换公开或私密。</p>
            </div>
          </article>
          <article class="step">
            <span class="step-num">4</span>
            <div>
              <h3>想分享给别人，就确认是公开</h3>
              <p>公开页面可以发链接给朋友；私密页面只适合自己用，不适合直接转发。</p>
            </div>
          </article>
        </div>
      </section>

      <section id="visibility" class="card">
        <h2>公开和私密有什么区别？</h2>
        <div class="compare-grid">
          <article class="compare-card">
            <h3>公开</h3>
            <p class="tagline">适合展示、分享、让别人访问。</p>
            <ul>
              <li>别人拿到链接后可以直接打开。</li>
              <li>可能会出现在「展示」页面里。</li>
              <li>适合小游戏、作品集、工具页、活动页。</li>
            </ul>
          </article>
          <article class="compare-card private">
            <h3>私密</h3>
            <p class="tagline">适合草稿、个人工具、暂时不想公开的页面。</p>
            <ul>
              <li>只有你登录后，从「管理」里打开才方便访问。</li>
              <li>别人直接打开链接会看到未公开提示。</li>
              <li>不会作为公开作品展示给其他用户。</li>
            </ul>
          </article>
        </div>
        <div class="note-box">
          <strong>简单记：</strong>
          想发给别人看，选公开；只是自己用，选私密。
        </div>
      </section>

      <section id="publish" class="card">
        <h2>发布模式是什么？</h2>
        <p class="lead">
          发布模式适合「我做了一个页面，希望大家打开时看到我准备好的内容或进度」这种场景。
        </p>
        <div class="example-list">
          <article>
            <h3>不开发布模式</h3>
            <p>每个用户打开页面后，页面里的数据更像是各用各的。比如待办清单、记账页面、个人备忘录，别人不应该看到你的内容。</p>
          </article>
          <article>
            <h3>开启发布模式</h3>
            <p>别人打开/刷新页面时，会直接使用你作为作者保存的那份内容。比如你做了一个官网、作品集、小游戏关卡，想让大家看到同一份初始内容。</p>
          </article>
        </div>
        <div class="note-box warning">
          <strong>建议：</strong>
          发布模式通常配合「公开」一起用。页面还是私密时，就算开启发布模式，别人也不一定能正常访问。
        </div>
      </section>

      <section id="problems" class="card">
        <h2>常见提示看不懂？</h2>
        <div class="faq-list">
          <article class="faq-item">
            <h3>看到「该页面未公开，请从控制台打开」怎么办？</h3>
            <p>意思是：这个页面现在不是公开页面，不能像普通网页一样直接访问。</p>
            <p>如果这是你自己的页面，请先回到 HtmlHub，登录后进入「管理」，点击这条记录旁边的「点击打开」。不要直接复制浏览器里的旧链接反复打开。</p>
            <p>如果这是别人发给你的链接，说明对方还没有把页面设为公开，或者页面暂时不能对外访问。请让作者在「管理」里切换为公开后再发给你。</p>
          </article>

          <article class="faq-item">
            <h3>我明明上传成功了，为什么别人打不开？</h3>
            <p>先检查「管理」里这条记录的可见性。如果显示「私密」，别人打不开是正常的；切换为公开后再分享。</p>
            <p>如果页面状态显示被拒绝，说明暂时不能公开展示，需要修改内容后再更新页面。</p>
          </article>

          <article class="faq-item">
            <h3>「展示」里为什么找不到我的页面？</h3>
            <p>通常是因为页面还是私密，或者内容没有通过检查。展示页主要放可以公开给大家看的页面。</p>
          </article>

          <article class="faq-item">
            <h3>更新 HTML 后，原来的链接会变吗？</h3>
            <p>不会。只要子域名不变，访问链接一般还是原来的那个。更新后刷新页面即可看到新内容。</p>
          </article>

          <article class="faq-item">
            <h3>页面里的数据会不会丢？</h3>
            <p>平台提供页面数据同步能力，但网站还在开发中，重要内容建议你在页面里手动导出或自己备份一份。</p>
          </article>

          <article class="faq-item">
            <h3>我应该选公开、私密，还是发布模式？</h3>
            <p>个人使用选私密；发给别人看选公开；想让别人看到你准备好的同一份页面内容，再考虑开启发布模式。</p>
          </article>
        </div>
      </section>

      <section class="card support-card">
        <h2>还是不会用？</h2>
        <p>可以联系站长反馈问题。描述时最好带上：你在哪个页面、点了哪个按钮、看到了什么提示。</p>
        <div class="contact-list">
          <div class="contact-item">
            <span>QQ 群：</span>
            <strong>1102036440</strong>
          </div>
          <div class="contact-item">
            <span>微信号：</span>
            <strong>shuijiaowang1</strong>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
.help-page {
  min-height: 100vh;
  background:
    radial-gradient(900px 420px at 20% -10%, rgb(var(--hh-brand-rgb) / 0.12), transparent 60%),
    radial-gradient(900px 420px at 100% 0%, rgb(var(--hh-brand-2-rgb) / 0.10), transparent 55%),
    var(--color-background);
}

.content {
  max-width: 1080px;
  margin: 0 auto;
  padding: 36px 24px;
}

.card {
  margin-bottom: 24px;
  padding: 26px;
  border: 1px solid var(--hh-border);
  border-radius: 16px;
  background: var(--hh-surface-solid);
  box-shadow: var(--hh-shadow-md);
}

.hero-card {
  text-align: center;
  padding: 46px 24px;
  background:
    radial-gradient(650px 280px at 20% 10%, rgb(var(--hh-brand-rgb) / 0.16), transparent 55%),
    radial-gradient(650px 280px at 90% 0%, rgb(var(--hh-brand-2-rgb) / 0.14), transparent 55%),
    linear-gradient(135deg, #f7fbff 0%, #f3f8ff 100%);
}

.eyebrow {
  margin: 0 0 8px;
  color: color-mix(in srgb, var(--hh-brand) 88%, #000 0%);
  font-size: 14px;
  font-weight: 800;
}

.hero-card h1 {
  margin: 0 0 10px;
  color: var(--color-heading);
  font-size: clamp(28px, 4.2vw, 42px);
  font-weight: 850;
  letter-spacing: -0.3px;
}

.subtitle {
  max-width: 680px;
  margin: 0 auto 20px;
  color: var(--hh-text-2);
  font-size: 15px;
  line-height: 1.7;
}

.quick-links {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 10px;
}

.quick-links a {
  padding: 8px 14px;
  border: 1px solid rgb(var(--hh-brand-rgb) / 0.16);
  border-radius: 999px;
  background: rgb(var(--hh-brand-rgb) / 0.10);
  color: color-mix(in srgb, var(--hh-brand) 88%, #000 0%);
  font-size: 13px;
  font-weight: 700;
}

.quick-links a:hover {
  text-decoration: none;
  background: rgb(var(--hh-brand-rgb) / 0.16);
}

h2 {
  margin: 0 0 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgb(var(--hh-brand-rgb) / 0.18);
  color: var(--color-heading);
  font-size: 22px;
  font-weight: 800;
  letter-spacing: -0.2px;
}

h3 {
  margin: 0 0 8px;
  color: var(--color-heading);
  font-size: 16px;
  font-weight: 800;
}

p {
  margin: 0;
  color: var(--hh-text-2);
  line-height: 1.7;
}

.steps,
.faq-list {
  display: grid;
  gap: 14px;
}

.step {
  display: flex;
  gap: 16px;
  padding: 16px;
  border: 1px solid rgba(2, 6, 23, 0.06);
  border-radius: 14px;
  background: color-mix(in srgb, var(--hh-surface-solid) 96%, #000 0%);
}

.step-num {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: color-mix(in srgb, var(--hh-brand) 86%, #000 0%);
  color: #fff;
  font-weight: 800;
}

.compare-grid,
.example-list {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

.compare-card,
.example-list article,
.faq-item {
  padding: 18px;
  border: 1px solid rgba(2, 6, 23, 0.06);
  border-radius: 14px;
  background:
    radial-gradient(420px 140px at 20% 0%, rgb(var(--hh-brand-rgb) / 0.10), transparent 60%),
    color-mix(in srgb, var(--hh-surface-solid) 95%, var(--hh-brand) 0%);
}

.compare-card.private {
  background:
    radial-gradient(420px 140px at 20% 0%, rgb(var(--hh-brand-2-rgb) / 0.10), transparent 60%),
    color-mix(in srgb, var(--hh-surface-solid) 95%, var(--hh-brand-2) 0%);
}

.tagline,
.lead {
  margin-bottom: 12px;
  color: var(--hh-text);
  font-weight: 600;
}

ul {
  margin: 12px 0 0;
  padding-left: 20px;
  color: var(--hh-text-2);
  line-height: 1.7;
}

.note-box {
  margin-top: 18px;
  padding: 14px 16px;
  border: 1px solid rgb(var(--hh-brand-rgb) / 0.14);
  border-radius: 12px;
  background: rgb(var(--hh-brand-rgb) / 0.06);
  color: var(--hh-text-2);
  line-height: 1.7;
}

.note-box strong {
  color: var(--color-heading);
}

.note-box.warning {
  border-color: rgba(245, 158, 11, 0.24);
  background: rgba(245, 158, 11, 0.08);
}

.faq-item p + p {
  margin-top: 10px;
}

.support-card {
  background:
    radial-gradient(650px 240px at 20% 0%, rgb(var(--hh-brand-rgb) / 0.10), transparent 55%),
    var(--hh-surface-solid);
}

.contact-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 18px;
}

.contact-item {
  display: inline-flex;
  gap: 8px;
  align-items: center;
  padding: 10px 12px;
  border: 1px solid rgba(2, 6, 23, 0.06);
  border-radius: 12px;
  background: color-mix(in srgb, var(--hh-surface-solid) 96%, #000 0%);
  color: var(--hh-text-2);
}

.contact-item strong {
  color: var(--color-heading);
}

@media (max-width: 760px) {
  .content {
    padding: 18px 12px;
  }

  .card {
    padding: 16px;
    border-radius: 14px;
  }

  .hero-card {
    padding: 30px 16px;
  }

  .compare-grid,
  .example-list {
    grid-template-columns: 1fr;
  }

  .step {
    gap: 12px;
    padding: 14px;
  }
}
</style>
