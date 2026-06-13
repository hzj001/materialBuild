<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getStats, type Stats } from "@/api/material";

defineOptions({ name: "Welcome" });

const stats = ref<Stats>({
  merchant_count: 0,
  order_count: 0,
  user_count: 0
});
const loading = ref(true);

onMounted(async () => {
  try {
    stats.value = await getStats();
  } finally {
    loading.value = false;
  }
});

const cards = [
  { key: "merchant_count", label: "入驻商家", icon: "ep:shop", color: "linear-gradient(135deg,#ff8c42,#ff5500)", desc: "覆盖三四线城市" },
  { key: "order_count", label: "平台订单", icon: "ep:document", color: "linear-gradient(135deg,#00b578,#08979c)", desc: "实时交易数据" },
  { key: "user_count", label: "注册用户", icon: "ep:user", color: "linear-gradient(135deg,#ffc300,#ff8f1f)", desc: "持续增长中" }
] as const;

const shortcuts = [
  { title: "城市管理", path: "/material/city", icon: "ep:office-building" },
  { title: "商家审核", path: "/material/merchant?status=0", icon: "ep:shop" },
  { title: "订单监控", path: "/material/order", icon: "ep:list" }
];
</script>

<template>
  <div v-loading="loading" class="welcome-page">
    <div class="hero">
      <div class="hero-content">
        <div class="hero-badge">materialBuild</div>
        <h1>建材通 · 运营指挥中心</h1>
        <p>像美团一样高效管理城市、商家与订单，赋能三四线城市建材零售</p>
      </div>
    </div>

    <el-row :gutter="16" class="stat-row">
      <el-col v-for="item in cards" :key="item.key" :xs="24" :sm="8">
        <div class="stat-card">
          <div class="stat-icon" :style="{ background: item.color }">
            <IconifyIconOffline :icon="item.icon" />
          </div>
          <div class="stat-body">
            <div class="stat-num">{{ stats[item.key] }}</div>
            <div class="stat-label">{{ item.label }}</div>
            <div class="stat-desc">{{ item.desc }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-card shadow="never" class="shortcut-card">
      <template #header><span class="card-title">快捷入口</span></template>
      <div class="shortcut-grid">
        <div
          v-for="s in shortcuts"
          :key="s.path"
          class="shortcut-item"
          @click="$router.push(s.path)"
        >
          <IconifyIconOffline :icon="s.icon" class="shortcut-icon" />
          <span>{{ s.title }}</span>
        </div>
      </div>
    </el-card>

    <el-card shadow="never" class="intro-card">
      <template #header><span class="card-title">平台能力</span></template>
      <el-row :gutter="12">
        <el-col :span="8"><el-tag effect="plain" type="warning">客户端 H5</el-tag></el-col>
        <el-col :span="8"><el-tag effect="plain" type="success">商家端 H5</el-tag></el-col>
        <el-col :span="8"><el-tag effect="plain">管理后台</el-tag></el-col>
      </el-row>
      <p class="intro-text">全链路覆盖选品、下单、配送、售后，Docker 容器组 materialBuild 一键部署。</p>
    </el-card>
  </div>
</template>

<style scoped>
.welcome-page { padding: 0 4px 16px; }
.hero {
  background: linear-gradient(135deg, #ff6b00 0%, #ff5500 50%, #e85d04 100%);
  border-radius: 16px; padding: 28px 24px; margin-bottom: 20px;
  color: #fff; position: relative; overflow: hidden;
}
.hero::after {
  content: ''; position: absolute; right: -30px; top: -30px;
  width: 160px; height: 160px; border-radius: 50%;
  background: rgba(255,255,255,0.1);
}
.hero-badge {
  display: inline-block; padding: 4px 10px; border-radius: 6px;
  background: rgba(255,255,255,0.2); font-size: 11px; font-weight: 600;
  letter-spacing: 1px; margin-bottom: 12px;
}
.hero h1 { font-size: 22px; font-weight: 800; margin: 0 0 8px; }
.hero p { font-size: 13px; opacity: 0.9; margin: 0; line-height: 1.6; }
.stat-row { margin-bottom: 16px; }
.stat-card {
  display: flex; align-items: center; gap: 16px;
  background: var(--el-bg-color); border-radius: 14px;
  padding: 20px; box-shadow: 0 2px 12px rgba(0,0,0,0.04);
  margin-bottom: 12px; border: 1px solid var(--el-border-color-lighter);
}
.stat-icon {
  width: 52px; height: 52px; border-radius: 14px;
  display: flex; align-items: center; justify-content: center;
  color: #fff; font-size: 24px; flex-shrink: 0;
}
.stat-num { font-size: 28px; font-weight: 800; color: var(--el-text-color-primary); line-height: 1.2; }
.stat-label { font-size: 14px; font-weight: 600; margin-top: 2px; }
.stat-desc { font-size: 12px; color: var(--el-text-color-secondary); margin-top: 2px; }
.card-title { font-weight: 700; font-size: 15px; }
.shortcut-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
.shortcut-item {
  display: flex; flex-direction: column; align-items: center; gap: 8px;
  padding: 20px 12px; border-radius: 12px; cursor: pointer;
  background: var(--el-fill-color-light); transition: all 0.2s;
}
.shortcut-item:hover { background: #fff8f0; color: #ff6b00; transform: translateY(-2px); }
.shortcut-icon { font-size: 28px; color: #ff6b00; }
.intro-text { margin: 16px 0 0; line-height: 1.8; color: var(--el-text-color-regular); font-size: 13px; }
</style>
