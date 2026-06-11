<template>
  <div class="home-page">
    <main class="content">
      <!-- 平台头部介绍 -->
      <section class="card hero-card">
        <h1>HtmlHub</h1>
        <p class="subtitle">一键上传 HTML，分配子域名，轻松托管静态页面</p>
        <p class="hero-desc">不用买服务器、不用搭数据库——让 AI 帮你写的网页，也能有专属链接、云端同步和分享能力。</p>
        <div class="hero-tags">
          <span>零部署</span>
          <span>多设备存储同步</span>
          <span>私密/公开访问</span>
          <span>一键分享</span>
        </div>
        <div class="hero-actions">
          <router-link to="/home/upload" class="hero-btn primary">立即上传</router-link>
          <router-link to="/home/manage" class="hero-btn secondary">管理页面</router-link>
          <router-link to="/home/showcase" class="hero-btn secondary">浏览作品</router-link>
        </div>
      </section>

      <!-- 核心功能 -->
      <section class="card">
        <h2>核心功能</h2>
        <div class="feature-grid">
          <div class="feature-item">
            <h4>简单上传</h4>
            <p>支持上传 HTML 文件或直接粘贴代码，无需复杂配置，上传即可访问</p>
          </div>
          <div class="feature-item">
            <h4>自定义子域名</h4>
            <p>可自定义子域名前缀，未填写时系统自动分配，生成专属访问链接</p>
          </div>
          <div class="feature-item">
            <h4>本地存储同步</h4>
            <p>内置 localStorage 云端同步脚本，支持数据导出、导入、上传与加载</p>
          </div>
          <div class="feature-item">
            <h4>权限控制</h4>
            <p>页面支持公开/私密切换，私密页面仅自己可访问，保护隐私</p>
          </div>
          <div class="feature-item">
            <h4>作品广场</h4>
            <p>浏览他人公开作品，支持点赞互动，发现有趣的 HTML 创意</p>
          </div>
          <div class="feature-item">
            <h4>浏览足迹</h4>
            <p>自动记录访问过的公开页面，方便回溯曾经看过的好作品</p>
          </div>
          <div class="feature-item">
            <h4>数据统计</h4>
            <p>统计页面访问次数与点赞数，清晰掌握页面使用情况</p>
          </div>
          <div class="feature-item">
            <h4>完整管理</h4>
            <p>集中管理所有页面，支持编辑、更新、删除与修改描述</p>
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
              <p>进入「上传」页面，填写子域名（可选），上传文件或粘贴完整 HTML 代码</p>
            </div>
          </div>
          <div class="step">
            <span class="step-num">2</span>
            <div class="step-content">
              <h4>获取访问链接</h4>
              <p>上传成功后系统自动分配子域名，生成形如 <code>xxx.htmlhub.cn</code> 的专属链接</p>
            </div>
          </div>
          <div class="step">
            <span class="step-num">3</span>
            <div class="step-content">
              <h4>管理与分享</h4>
              <p>在「管理」中查看、编辑页面，设置公开/私密后分享给他人；更多说明见「帮助」</p>
            </div>
          </div>
        </div>
        <p class="guide-footer">
          还没注册？先
          <router-link to="/home?auth=register">免费注册</router-link>
          ，即可开始上传你的第一个页面。
        </p>
      </section>

      <!-- AI 创作提示 -->
      <section class="card">
        <h2>AI 快速生成页面</h2>
        <p class="tips-title">复制下方提示词发给 AI，即可生成可直接托管的完整 HTML 页面：</p>
        <div class="code-block-wrap">
          <pre class="code-block">{{ aiPrompt }}</pre>
          <button type="button" class="copy-btn" @click="copyPrompt">复制提示词</button>
        </div>
        <p class="tips-desc">生成完成后，复制代码到「上传」页面即可使用。需要 localStorage 同步的页面，可让 AI 参考平台帮助文档。</p>
      </section>

      <!-- 联系与支持 -->
      <section class="card">
        <h2>联系站长</h2>
        <p class="platform-desc">
          本网站主要为非程序员打造，不用自己部署服务器、无需搭建数据库，解决纯静态网页无法数据同步、无法分享页面的问题。平台目前尚处于开发阶段，请注意备份自己的数据。<br>
          ✅ 如需纯净子域名、修改使用上限，可联系站长；<br>
          ✅ 如遇到网站 bug、功能建议，或是 VibeCoding 创作过程中遇到问题，均可联系站长帮忙解答。
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
            <span>小红书：</span>
            <a target="_blank" rel="noopener noreferrer" href="https://www.xiaohongshu.com/discovery/item/69fde400000000003700d847?source=webshare&xhsshare=pc_web&xsec_token=AB7ubb39Ok16_hMhgDMM0uEnkK2cI7qylkRrdSFubcxXM=&xsec_source=pc_share">睡觉王</a>
          </div>
          <div class="contact-item">
            <span>独立部署：</span>
            <span>支持个人网站 / 服务器 / 数据库定制部署服务</span>
          </div>
        </div>
      </section>

      <!-- 赏赞支持 -->
      <section class="card donate-card">
        <h2>支持站长</h2>
        <p class="donate-desc">
          HtmlHub 由个人开发与维护，服务器与域名费用均由站长承担。如果你觉得平台有帮助，欢迎扫码赏赞支持，每一份鼓励都是继续优化的动力。
        </p>
        <div class="donate-grid">
          <figure class="donate-item">
            <img src="/donate/wechat.png" alt="微信赞赏码" width="180" height="180" loading="lazy" />
            <figcaption>微信赞赏</figcaption>
          </figure>
          <figure class="donate-item">
            <img src="/donate/alipay.png" alt="支付宝收款码" width="180" height="180" loading="lazy" />
            <figcaption>支付宝</figcaption>
          </figure>
        </div>
        <p class="donate-note">赏赞完全自愿，不影响任何功能使用。也可以帮忙分享给有需要的朋友。</p>
      </section>

      <!-- 开发日志 -->
      <section class="card changelog-card">
        <h2>开发日志</h2>
        <p class="changelog-desc">记录平台主要更新，按天汇总。点击日期可展开或折叠。</p>
        <div class="changelog-list">
          <div
            v-for="group in changelog"
            :key="group.date"
            class="changelog-group"
            :class="{ expanded: expandedDates.has(group.date) }"
          >
            <button
              type="button"
              class="changelog-header"
              :aria-expanded="expandedDates.has(group.date)"
              @click="toggleDate(group.date)"
            >
              <span class="changelog-date">{{ group.label }}</span>
              <span class="changelog-meta">{{ group.items.length }} 项更新</span>
              <span class="changelog-chevron" aria-hidden="true">{{ expandedDates.has(group.date) ? '−' : '+' }}</span>
            </button>
            <ul v-show="expandedDates.has(group.date)" class="changelog-items">
              <li v-for="(item, index) in group.items" :key="index">{{ item }}</li>
            </ul>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

interface ChangelogGroup {
  date: string
  label: string
  items: string[]
}

const changelog: ChangelogGroup[] = [
  {
    date: '2026-06-11',
    label: '2026年6月11日',
    items: [
      '上传页支持拖拽文件，拖入即可上传',
      '管理页新增回收站，误删可恢复，也可彻底删除',
      '新增个人中心与浏览足迹，可查看和修改部分账号数据',
      '管理端支持调整用户配额，如 HTML 数量、累计上传上限等',
    ],
  },
  {
    date: '2026-05-19',
    label: '2026年5月19日',
    items: [
      '全站界面小幅美化，阅读与操作更清晰',
      '新增帮助文档页，常见问题集中解答',
      '新增发布者模式，方便在页面内展示作者信息',
      '管理页补充点赞统计展示',
      '作品点赞功能上线，含点赞记录与相关接口',
    ],
  },
  {
    date: '2026-05-09',
    label: '2026年5月9日',
    items: [
      '上传编辑区新增复制、粘贴、清空按钮，手机端更好用',
    ],
  },
  {
    date: '2026-05-08',
    label: '2026年5月8日',
    items: [
      '修复并优化 localStorage 云端同步脚本',
      '多处细节适配手机端显示与操作',
    ],
  },
  {
    date: '2026-05-07',
    label: '2026年5月7日',
    items: [
      '大幅优化各页面布局与视觉样式',
      '同步脚本支持跳转主页后自动同步数据',
      '调整默认公开策略，审核中页面也可分享预览',
      '优化登录注册流程，未登录也可浏览首页与公开内容',
      '补充无域名场景下的路由测试',
    ],
  },
  {
    date: '2026-05-04',
    label: '2026年5月4日',
    items: [
      '未登录用户也可访问首页',
      '前端代码内嵌编译，修复线上资源加载失败',
      '补充 SEO 与站点图标',
      '同步脚本精简体积并优化提示文案',
      '更新部署与运行配置',
    ],
  },
  {
    date: '2026-05-03',
    label: '2026年5月3日',
    items: [
      '拆分同步脚本，结构更清晰',
      '登录注册增加频率限制与 IP 限流',
      '新增 HTML 更新逻辑与上传校验',
      '设定默认配额：单页 ≤ 1MB、同步数据 ≤ 100KB、在架 ≤ 20 页、累计 ≤ 100 页',
    ],
  },
  {
    date: '2026-05-02',
    label: '2026年5月2日',
    items: [
      '新增管理后台，支持用户管理与 HTML 管理',
      '新增页面访问次数统计与展示',
    ],
  },
  {
    date: '2026-05-01',
    label: '2026年5月1日',
    items: [
      '未公开页面仅自己可见，支持删除与公开/私密切换',
      '前端页面拆分，便于后续维护扩展',
      '自定义子域名时自动拼接随机后缀，避免冲突',
      '调整数据模型与路由，保留四级域名访问能力',
    ],
  },
  {
    date: '2026-04-18',
    label: '2026年4月18日',
    items: [
      '项目起步，搭建基础上传托管与无域名路由能力',
    ],
  },
]

const expandedDates = ref(new Set([changelog[0].date]))

const toggleDate = (date: string) => {
  const next = new Set(expandedDates.value)
  if (next.has(date)) {
    next.delete(date)
  } else {
    next.add(date)
  }
  expandedDates.value = next
}

const aiPrompt = `请用 HTML + CSS + JS 生成一个完整的单文件页面，要求：
1. 所有代码写在一个 HTML 文件里，可直接保存上传
2. 如需保存用户数据，使用 localStorage
3. 页面需求：【在这里描述你的功能，例如：待办清单、记账本、小游戏】
4. 风格简洁美观，电脑端/移动端友好`

const copyPrompt = async () => {
  try {
    await navigator.clipboard.writeText(aiPrompt)
    ElMessage.success('提示词已复制')
  } catch {
    ElMessage.error('复制失败，请手动选择复制')
  }
}
</script>

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

.card {
  background: var(--hh-surface-solid);
  border: 1px solid var(--hh-border);
  border-radius: 16px;
  padding: 26px;
  margin-bottom: 24px;
  box-shadow: var(--hh-shadow-md);
}

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
  margin-bottom: 10px;
}

.hero-desc {
  font-size: 14px;
  color: var(--hh-text-3);
  max-width: 560px;
  margin: 0 auto 20px;
  line-height: 1.6;
}

.hero-tags {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 24px;
}

.hero-tags span {
  padding: 7px 14px;
  background: rgb(var(--hh-brand-rgb) / 0.10);
  color: color-mix(in srgb, var(--hh-brand) 88%, #000 0%);
  border-radius: 20px;
  font-size: 13px;
  border: 1px solid rgb(var(--hh-brand-rgb) / 0.16);
}

.hero-actions {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 12px;
}

.hero-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 10px 22px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 700;
  text-decoration: none;
  transition: background 0.15s, transform 0.1s;
}

.hero-btn.primary {
  background: color-mix(in srgb, var(--hh-brand) 86%, #000 0%);
  color: #fff;
  border: 1px solid transparent;
}

.hero-btn.primary:hover {
  background: color-mix(in srgb, var(--hh-brand) 78%, #000 0%);
  text-decoration: none;
}

.hero-btn.secondary {
  background: var(--hh-surface-solid);
  color: var(--color-heading);
  border: 1px solid var(--hh-border);
}

.hero-btn.secondary:hover {
  background: rgb(var(--hh-brand-rgb) / 0.08);
  text-decoration: none;
}

.hero-btn:active {
  transform: scale(0.98);
}

.card h2 {
  font-size: 22px;
  font-weight: 800;
  color: var(--color-heading);
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgb(var(--hh-brand-rgb) / 0.18);
  letter-spacing: -0.2px;
}

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

.guide-steps {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.step {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 14px;
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

.step-content code {
  padding: 2px 6px;
  border-radius: 6px;
  background: rgba(2, 6, 23, 0.05);
  font-size: 13px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.guide-footer {
  margin: 18px 0 0;
  font-size: 14px;
  color: var(--hh-text-3);
  text-align: center;
}

.tips-title {
  font-size: 15px;
  color: var(--hh-text-2);
  margin-bottom: 12px;
}

.code-block-wrap {
  position: relative;
  margin-bottom: 8px;
}

.code-block {
  padding: 16px 16px 48px;
  background: rgba(2, 6, 23, 0.04);
  border-radius: 12px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  font-size: 13px;
  color: var(--hh-text);
  border: 1px solid rgba(2, 6, 23, 0.06);
  line-height: 1.6;
  overflow: auto;
  white-space: pre-wrap;
  margin: 0;
}

.copy-btn {
  position: absolute;
  right: 10px;
  bottom: 10px;
  padding: 6px 12px;
  border-radius: 8px;
  border: 1px solid var(--hh-border);
  background: var(--hh-surface-solid);
  color: var(--color-heading);
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
}

.copy-btn:hover {
  background: rgb(var(--hh-brand-rgb) / 0.10);
}

.tips-desc {
  font-size: 14px;
  color: var(--hh-text-3);
  margin: 0;
}

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

.contact-item span:last-child,
.contact-item a {
  color: var(--hh-text-2);
}

.donate-desc {
  font-size: 15px;
  color: var(--hh-text-2);
  line-height: 1.6;
  margin: 0 0 24px;
}

.donate-grid {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 32px;
}

.donate-item {
  margin: 0;
  text-align: center;
}

.donate-item img {
  display: block;
  width: 180px;
  height: 180px;
  border-radius: 14px;
  border: 1px solid var(--hh-border);
  background: #fff;
  box-shadow: 0 8px 20px rgba(2, 6, 23, 0.08);
}

.donate-item figcaption {
  margin-top: 10px;
  font-size: 14px;
  font-weight: 700;
  color: var(--color-heading);
}

.donate-note {
  margin: 20px 0 0;
  font-size: 13px;
  color: var(--hh-text-3);
  text-align: center;
}

.changelog-desc {
  font-size: 14px;
  color: var(--hh-text-3);
  margin: -8px 0 16px;
}

.changelog-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.changelog-group {
  border: 1px solid rgba(2, 6, 23, 0.06);
  border-radius: 12px;
  overflow: hidden;
  background: color-mix(in srgb, var(--hh-surface-solid) 98%, #000 0%);
}

.changelog-group.expanded {
  border-color: rgb(var(--hh-brand-rgb) / 0.22);
}

.changelog-header {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 12px 14px;
  border: none;
  background: transparent;
  cursor: pointer;
  text-align: left;
  font: inherit;
  color: inherit;
}

.changelog-header:hover {
  background: rgb(var(--hh-brand-rgb) / 0.06);
}

.changelog-date {
  font-size: 15px;
  font-weight: 800;
  color: var(--color-heading);
}

.changelog-meta {
  font-size: 13px;
  color: var(--hh-text-3);
}

.changelog-chevron {
  margin-left: auto;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  background: rgb(var(--hh-brand-rgb) / 0.10);
  color: var(--color-heading);
  font-size: 16px;
  font-weight: 700;
  flex-shrink: 0;
}

.changelog-items {
  margin: 0;
  padding: 4px 14px 14px 34px;
  list-style: disc;
}

.changelog-items li {
  font-size: 14px;
  color: var(--hh-text-2);
  line-height: 1.55;
  margin-bottom: 6px;
}

.changelog-items li:last-child {
  margin-bottom: 0;
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

  .donate-item img {
    width: 150px;
    height: 150px;
  }

  .changelog-header {
    padding: 10px 12px;
    flex-wrap: wrap;
  }

  .changelog-chevron {
    margin-left: 0;
  }
}

@media (max-width: 1024px) and (min-width: 641px) {
  .feature-item {
    grid-column: span 6;
  }
}
</style>
