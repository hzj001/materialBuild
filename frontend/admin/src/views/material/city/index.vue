<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { message } from "@/utils/message";
import {
  getCities,
  createCity,
  updateCityStatus,
  type City
} from "@/api/material";

defineOptions({ name: "MaterialCity" });

const loading = ref(false);
const list = ref<City[]>([]);
const dialogVisible = ref(false);
const form = reactive({ name: "", province: "", code: "" });

const loadData = async () => {
  loading.value = true;
  try {
    list.value = await getCities();
  } finally {
    loading.value = false;
  }
};

const handleAdd = async () => {
  if (!form.name || !form.code) {
    message("请填写城市名称和编码", { type: "warning" });
    return;
  }
  await createCity({ ...form });
  message("添加成功", { type: "success" });
  dialogVisible.value = false;
  form.name = "";
  form.province = "";
  form.code = "";
  loadData();
};

const toggleStatus = async (row: City) => {
  const status = row.status === 1 ? 0 : 1;
  await updateCityStatus(row.id, status);
  message("状态已更新", { type: "success" });
  loadData();
};

onMounted(loadData);
</script>

<template>
  <div>
    <el-card shadow="never">
      <div class="toolbar">
        <span class="title">服务城市列表</span>
        <el-button type="primary" @click="dialogVisible = true">添加城市</el-button>
      </div>
      <el-table v-loading="loading" :data="list" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="城市" min-width="120" />
        <el-table-column prop="province" label="省份" min-width="100" />
        <el-table-column prop="code" label="编码" min-width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? "启用" : "禁用" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="toggleStatus(row)">
              {{ row.status === 1 ? "禁用" : "启用" }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="添加城市" width="420px">
      <el-form label-width="80px">
        <el-form-item label="城市名称" required>
          <el-input v-model="form.name" placeholder="如：遵义" />
        </el-form-item>
        <el-form-item label="省份">
          <el-input v-model="form.province" placeholder="如：贵州" />
        </el-form-item>
        <el-form-item label="城市编码" required>
          <el-input v-model="form.code" placeholder="如：520300" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAdd">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.title {
  font-size: 16px;
  font-weight: 600;
}
</style>
