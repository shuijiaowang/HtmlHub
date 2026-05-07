<template>
  <div class="home-page">
    <main class="content">
      <section class="card">
        <h1>展示页</h1>
        <p class="subtitle">这里将用于展示大家都上传了什么内容，后续开发中。</p>
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

.card {
  background: var(--hh-surface-solid);
  border: 1px solid var(--hh-border);
  border-radius: var(--hh-radius-md);
  padding: 18px;
  margin-bottom: 18px;
  box-shadow: var(--hh-shadow-md);
}

.subtitle {
  color: #666;
  margin-bottom: 20px;
}
</style>
