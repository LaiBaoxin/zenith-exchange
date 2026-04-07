import axios, { type AxiosInstance, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios';

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
        const token = localStorage.getItem('token');
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
            console.log("res.data", res.data)
            return res.data;
        }

        const errorMsg = res.msg || '业务请求失败';
        console.error('业务错误:', errorMsg);
        return Promise.reject(new Error(errorMsg));
    },
    (error) => {
        const message = error.response?.data?.msg || error.message || '网络连接异常';

        if (error.response?.status === 401) {
            console.warn('登录已过期，请重新登录');
            window.localStorage.clear()
        }

        return Promise.reject(new Error(message));
    }
);

export default request;
