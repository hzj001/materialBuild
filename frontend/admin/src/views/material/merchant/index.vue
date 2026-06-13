<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { message } from "@/utils/message";
import {
  getMerchants,
  updateMerchantStatus,
  type Merchant
} from "@/api/material";

defineOptions({ name: "MaterialMerchant" });

const statusMap: Record<number, { text: string; type: string }> = {
  0: { text: "待审核", type: "warning" },
  1: { text: "营业中", type: "success" },
  2: { text: "休息中", type: "info" },
  3: { text: "已禁用", type: "danger" }
};

const loading = ref(false);
const list = ref<Merchant[]>([]);
const total = ref(0);
const query = reactive({ status: "", page: 1, page_size: 20 });

const loadData = async () => {
  loading.value = true;
  try {
    const params: Record<string, any> = {
      page: query.page,
      page_size: query.page_size
    };
    if (query.status !== "") params.status = query.status;
    const data = await getMerchants(params);
    list.value = data.list || [];
    total.value = data.total || 0;
  } finally {
    loading.value = false;
  }
};

const setStatus = async (row: Merchant, status: number) => {
  await updateMerchantStatus(row.id, status);
  message("操作成功", { type: "success" });
  loadData();
};

onMounted(loadData);
</script>

<template>
  <div>
    <el-card shadow="never">
      <div class="toolbar">
        <el-radio-group v-model="query.status" @change="loadData">
          <el-radio-button label="">全部</el-radio-button>
          <el-radio-button label="0">待审核</el-radio-button>
          <el-radio-button label="1">营业中</el-radio-button>
          <el-radio-button label="3">已禁用</el-radio-button>
        </el-radio-group>
      </div>
      <el-table v-loading="loading" :data="list" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="店铺名称" min-width="140" />
        <el-table-column prop="contact_phone" label="联系电话" width="130" />
        <el-table-column prop="address" label="地址" min-width="180" show-overflow-tooltip />
        <el-table-column prop="rating" label="评分" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="(statusMap[row.status]?.type as any) || 'info'">
              {{ statusMap[row.status]?.text || "未知" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <template v-if="row.status === 0">
              <el-button link type="primary" @click="setStatus(row, 1)">通过</el-button>
              <el-button link type="danger" @click="setStatus(row, 3)">拒绝</el-button>
            </template>
            <el-button
              v-else-if="row.status === 1"
              link
              type="danger"
              @click="setStatus(row, 3)"
            >禁用</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="query.page"
        class="mt-4"
        layout="total, prev, pager, next"
        :total="total"
        @current-change="loadData"
      />
    </el-card>
  </div>
</template>

<style scoped>
.toolbar {
  margin-bottom: 16px;
}
.mt-4 {
  margin-top: 16px;
  justify-content: flex-end;
}
</style>
