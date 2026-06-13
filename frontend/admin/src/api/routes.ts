/** 使用静态路由，不再从后端拉取动态菜单 */
export const getAsyncRoutes = () => {
  return Promise.resolve({ success: true, data: [] as Array<any> });
};
