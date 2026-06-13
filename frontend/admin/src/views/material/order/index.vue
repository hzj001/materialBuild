<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { getOrders, type Order } from "@/api/material";

defineOptions({ name: "MaterialOrder" });

const statusMap: Record<number, { text: string; type: string }> = {
  0: { text: "待支付", type: "warning" },
  1: { text: "已支付", type: "success" },
  2: { text: "配送中", type: "primary" },
  3: { text: "已完成", type: "success" },
  4: { text: "已取消", type: "info" },
  5: { text: "售后中", type: "danger" }
};

const loading = ref(false);
const list = ref<Order[]>([]);
const total = ref(0);
const query = reactive({ page: 1, page_size: 20 });

const loadData = async () => {
  loading.value = true;
  try {
    const data = await getOrders(query);
    list.value = data.list || [];
    total.value = data.total || 0;
  } finally {
    loading.value = false;
  }
};

onMounted(loadData);
</script>

<template>
  <div>
    <el-card shadow="never">
      <el-table v-loading="loading" :data="list" stripe>
        <el-table-column prop="order_no" label="订单号" min-width="180" />
        <el-table-column label="商品" min-width="200">
          <template #default="{ row }">
            <div v-for="(item, i) in row.items || []" :key="i" class="item-line">
              {{ item.product_name }} × {{ item.quantity }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="实付" width="100">
          <template #default="{ row }">
            <span class="amount">¥{{ row.pay_amount }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="(statusMap[row.status]?.type as any) || 'info'">
              {{ statusMap[row.status]?.text || "未知" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="下单时间" min-width="170" />
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
.item-line {
  font-size: 13px;
  line-height: 1.6;
}
.amount {
  color: #e85d04;
  font-weight: 600;
}
.mt-4 {
  margin-top: 16px;
  justify-content: flex-end;
}
</style>
