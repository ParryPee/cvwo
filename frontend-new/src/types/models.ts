interface User {
    id: number;
    username: string;
    created_at: string;
}

interface Topic {
    id: number;
    title: string;
    description: string;
    created_at: string;
    user_id: number;
    username: string;
    post_count: number;
}

interface Post {
    id: number;
    title: string;
    content: string;
    likes: number;
    created_at: string;
    updated_at: string;
    topic_id: number;
    user_id: number;
    topic_title: string;
    username: string;
    liked_by_user: boolean;
}

interface Comment {
    id: number;
    content: string;
    likes: number;
    created_at: string;
    updated_at: string;
    post_id: number;
    user_id: number;
    deleted: boolean;
    parent_comment_id: {Int64: number, Valid: boolean};
    liked_by_user: boolean;
    username: string;
}
interface SearchResult {
    posts: Post[];
    topics: Topic[];
}

export type { User, Topic, Post, Comment, SearchResult };