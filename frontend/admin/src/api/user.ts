import { http } from "@/utils/http";

export type UserResult = {
  success: boolean;
  data: {
    avatar: string;
    username: string;
    nickname: string;
    roles: Array<string>;
    permissions: Array<string>;
    accessToken: string;
    refreshToken: string;
    expires: Date;
  };
};

export type RefreshTokenResult = {
  success: boolean;
  data: {
    accessToken: string;
    refreshToken: string;
    expires: Date;
  };
};

type LoginPayload = {
  username?: string;
  password?: string;
};

type BackendLoginResult = {
  token: string;
  user: {
    phone: string;
    nickname: string;
    avatar?: string;
    role: string;
  };
};

/** 登录 — 对接建材通 Go 后端 */
export const getLogin = async (data?: object) => {
  const body = data as LoginPayload;
  const result = await http.request<BackendLoginResult>(
    "post",
    "/api/v1/auth/login",
    {
      data: {
        phone: body?.username,
        password: body?.password,
        role: "admin"
      }
    }
  );

  const expires = new Date();
  expires.setDate(expires.getDate() + 7);

  return {
    success: true,
    data: {
      avatar: result?.user?.avatar || "",
      username: result?.user?.phone || body?.username,
      nickname: result?.user?.nickname || "管理员",
      roles: ["admin"],
      permissions: ["*:*:*"],
      accessToken: result?.token,
      refreshToken: result?.token,
      expires
    }
  } as UserResult;
};

/** 刷新 token（后端暂未实现，复用当前 token） */
export const refreshTokenApi = async (data?: { refreshToken?: string }) => {
  const expires = new Date();
  expires.setDate(expires.getDate() + 7);
  return {
    success: true,
    data: {
      accessToken: data?.refreshToken || "",
      refreshToken: data?.refreshToken || "",
      expires
    }
  } as RefreshTokenResult;
};
