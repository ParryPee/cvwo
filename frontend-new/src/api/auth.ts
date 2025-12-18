import client from './client';
import type { User } from '../types/models';


export const login = async (username: string): Promise<User> => {
    const response = await client.post<User>('users/login', {username});
    return response.data;
}
export const register = async (username: string): Promise<User> => {
    const response = await client.post<User>('users/register', {username});
    return response.data;
}
export const fetchCurrentUser = async (): Promise<User> => {
    const response = await client.get<User>('users/me');
    if (response.status !== 200) {
        return Promise.reject('No authenticated user');
    }
    return response.data;
}
export const logout = async (): Promise<void> => {
    const resp = await client.post('users/logout');
    if (resp.status !== 200) {
        throw new Error('Logout failed');
    }
}