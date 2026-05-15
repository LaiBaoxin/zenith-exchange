import axios, { type AxiosInstance, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios';

// 定义后端通用的返回结构
interface BackEndResponse<T = any> {
    code: number;
    data: T;
    msg: string;
}

const request: AxiosInstance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
});

// 请求拦截器
request.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        const token = window.localStorage.getItem('zenith_auth_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

// 响应拦截器
request.interceptors.response.use(
    (response: AxiosResponse<BackEndResponse>) => {
        const res = response.data;
        if (res.code === 200) {
            return res.data;
        }

        const errorMsg = res.msg || '业务请求失败';
        return Promise.reject(new Error(errorMsg));
    },
    (error) => {
        const message = error.response?.data?.msg || error.message || '网络连接异常';

        if (error.response?.status === 401) {
            if (error.config.url?.includes('/auth/login')) {
                return Promise.reject(new Error(message));
            }

            console.error('鉴权失效:', error.config.url);
            localStorage.removeItem('zenith_auth_token');
            localStorage.removeItem('user_address');
            window.dispatchEvent(new CustomEvent('auth:expired'));
        }

        return Promise.reject(new Error(message));
    }
);

export default request;
