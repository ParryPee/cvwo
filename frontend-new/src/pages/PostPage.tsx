import { use, useEffect, useState } from "react";
import { useNavigate, useParams, Link } from "react-router-dom";
import { fetchPostById, fetchCommentsByPostId } from "../api/forum";
import type { Post, Comment } from "../types/models";
import { useAuth } from "../context/AuthContext";
import { timeAgo, formatDate } from "../utils/date";
import Box from "@mui/material/Box";
import {
	Container,
	CircularProgress,
	Typography,
	Grid,
	CardActionArea,
	Alert,
	Button,
} from "@mui/material";
const TopicPage = () => {
	const navigate = useNavigate();
	const { topicId } = useParams<{ topicId: string }>();
	const { postId } = useParams<{ postId: string }>();
	const [loading, setLoading] = useState(true);
	const [post, setPost] = useState<Post | null>(null);
	const [comments, setComments] = useState<Comment[]>([]);
	const [error, setError] = useState<string | null>(null);
	const { isAuthenticated } = useAuth();
	useEffect(() => {
		const loadPostAndComments = async () => {
			if (!postId || !topicId) return;
			try {
				setLoading(true);
				const postIdNum = parseInt(postId, 10);
				if (isNaN(postIdNum)) {
					setError("Invalid Post ID");
					setLoading(false);
					return;
				}
				const postData = await fetchPostById(postIdNum);
				const commentsData = await fetchCommentsByPostId(postIdNum);
				setPost(postData);
				setComments(commentsData);
			} catch (error) {
				console.error("Failed to fetch topic or posts:", error);
				setError("Failed to load topic. It may not exist.");
			} finally {
				setLoading(false);
			}
		};
		loadPostAndComments();
	}, [postId, topicId]);
	if (loading) {
		return (
			<Box display="flex" justifyContent="center" mt={4}>
				<CircularProgress />
			</Box>
		);
	}
	if (error || !post) {
		return (
			<Container>
				<Alert severity="error">{error || "Post not found."}</Alert>
			</Container>
		);
	}
	return <Container sx={{ mt: 4 }}>Post Page</Container>;
};

export default TopicPage;
