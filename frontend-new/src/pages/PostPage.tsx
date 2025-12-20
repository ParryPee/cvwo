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
	return (
		<Container
			sx={{
				mt: 4,
				padding: "16px",
			}}
		>
			<Box sx={{ border: "1px solid #ccc", borderRadius: "8px", p: 2 }}>
				<Typography variant="h4" gutterBottom>
					{post.title}
				</Typography>
				<Typography variant="caption" color="text.secondary">
					Posted by • {post.created_by_username} •{" "}
					{timeAgo(post.created_at)} (on {formatDate(post.created_at)}
					)
				</Typography>
				<Typography
					variant="subtitle1"
					color="text.secondary"
					gutterBottom
				>
					{post.content}
				</Typography>
			</Box>

			<Box mt={4}>
				<Box
					sx={{
						display: "flex",
						justifyContent: "space-between",
						alignItems: "center",
						mb: 2,
					}}
				>
					<Typography variant="h6">Comments</Typography>
					{isAuthenticated && (
						<Button
							variant="contained"
							color="primary"
							onClick={() =>
								navigate(
									`/topics/${topicId}/posts/${postId}/comments/create`
								)
							}
						>
							Create Comment
						</Button>
					)}
				</Box>
				{!comments || comments.length === 0 ? (
					<Typography variant="body2" color="text.secondary" mt={2}>
						No comments yet. Be the first to comment!
					</Typography>
				) : (
					comments.map((comment) => (
						<Grid key={comment.id} size={{ xs: 12, sm: 6 }}>
							<CardActionArea
								sx={{
									border: "1px solid #ddd",
									borderRadius: "8px",
									padding: "16px",
									marginBottom: "16px",
								}}
								//Todo: Link to comment detail page, subreplies etc
							>
								<Typography
									variant="body1"
									sx={{
										display: "-webkit-box",
										WebkitLineClamp: 2,
										WebkitBoxOrient: "vertical",
										overflow: "hidden",
										textOverflow: "ellipsis",
									}}
								>
									{comment.content}
								</Typography>
								<Typography
									variant="caption"
									color="text.secondary"
								>
									Posted by • {comment.created_by_username} •{" "}
									{timeAgo(comment.created_at)} (on{" "}
									{formatDate(comment.created_at)})
								</Typography>
								<Box my={2} />
							</CardActionArea>
						</Grid>
					))
				)}
			</Box>
		</Container>
	);
};

export default TopicPage;
