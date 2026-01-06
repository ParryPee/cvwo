import client from './client';
import type { Post, Topic, Comment } from '../types/models';


export const fetchAllTopics = async (): Promise<Topic[]> => {
    const response = await client.get<Topic[]>('topics');
    return response.data;
}

export const fetchTopicById = async (id: number): Promise<Topic> => {
    const response = await client.get<Topic>(`topics/${id}`);
    return response.data;
}
export const fetchPostsByTopicId = async (topicId: number): Promise<Post[]> => {
    const response = await client.get<Post[]>(`topics/${topicId}/posts`);
    return response.data;
}
export const createPost = async (postData: { topic_id: number; content: string; title: string; user_id: number }): Promise<Post> => {
    const response = await client.post<Post>('posts', postData);
    return response.data;
}
export const createTopic = async (topicData: { title: string; description: string; created_by: number }): Promise<Topic> => {
    const response = await client.post<Topic>('topics', topicData);
    return response.data;
}
export const fetchPostById = async (postId: number): Promise<Post> => {
    const response = await client.get<Post>(`posts/${postId}`);
    return response.data;
}
export const fetchCommentsByPostId = async (postId: number): Promise<Comment[]> => {
    const response = await client.get<Comment[]>(`posts/${postId}/comments`);
    return response.data;
}
export const createComment = async (commentData: { post_id: number; content: string; user_id: number; parent_id?: number | null }): Promise<Comment> => {
    const response = await client.post<Comment>('comments', commentData);
    return response.data;
}
export const fetchCommentById = async (commentId: number): Promise<Comment> => {
    const response = await client.get<Comment>(`comments/${commentId}`);
    return response.data;
}
export const likeComment = async (commentId: number): Promise<void> => {
    await client.post(`comments/${commentId}/like`);
}
export const likePost = async (postId: number): Promise<void> =>{
    await client.post(`posts/${postId}/like`);
}
export const deletePost = async (postID: number): Promise<void> =>{
    await client.delete(`posts/${postID}`)
}
export const searchPost = async (query: string): Promise<Post[]> =>{
    const response  = await client.get(`search?q=${query}`)
    return response.data
}
export const deleteComment = async(commentID: number): Promise<void> =>{
    await client.delete(`comments/${commentID}`)
}

