interface User {
    id: number;
    username: string;
    created_at: string;
}

interface Topic {
    id: number;
    title: string;
    description: string;
    created_by: number;
    created_at: string;
    created_by_username: string;
}

interface Post {
    id: number;
    title: string;
    content: string;
    topic_id: number;
    user_id: number;
    created_at: string;
    updated_at: string;
    created_by_username: string;
}

interface Comment {
    id: number;
    content: string;
    post_id: number;
    user_id: number;
    likes: number;
    parent_id: number | null;
    created_at: string;
    updated_at: string;
    created_by_username: string;
    liked_by_user: boolean;
}

export type { User, Topic, Post, Comment };