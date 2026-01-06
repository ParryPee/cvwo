import { use, useEffect, useState } from "react";
import { useNavigate, useParams, Link } from "react-router-dom";
import {
	fetchPostById,
	fetchCommentsByPostId,
	createComment,
	likeComment,
	likePost,
	deletePost,
	deleteComment,
} from "../api/forum";
import type { Post, Comment } from "../types/models";
import { useAuth } from "../context/AuthContext";
import { timeAgo, formatDate } from "../utils/date";
import Box from "@mui/material/Box";
import TextBox from "../components/TextBox";
import FavoriteBorderIcon from "@mui/icons-material/FavoriteBorder";
import IconButton from "@mui/material/IconButton";
import FavoriteIcon from "@mui/icons-material/Favorite";
import DeleteIcon from "@mui/icons-material/Delete";

import {
	Container,
	CircularProgress,
	Typography,
	Grid,
	Card,
	Alert,
	Button,
	Icon,
} from "@mui/material";
const PostPage = () => {
	const navigate = useNavigate();
	const { topicId } = useParams<{ topicId: string }>();
	const { postId } = useParams<{ postId: string }>();
	const [loading, setLoading] = useState(true);
	const [post, setPost] = useState<Post | null>(null);
	const [comments, setComments] = useState<Comment[]>([]);
	const [error, setError] = useState<string | null>(null);
	const { isAuthenticated, user } = useAuth();

	const handleSubmit = async (
		text: string,
		parentId: number | null = null
	) => {
		if (!isAuthenticated || !postId) {
			navigate("/login");
			return;
		}
		try {
			const postIdNum = parseInt(postId, 10);
			if (isNaN(postIdNum)) {
				setError("Invalid Post ID");
				return;
			}
			await createComment({
				post_id: postIdNum,
				content: text,
				user_id: user!.id,
				parent_id: parentId,
			});
			navigate(0);
		} catch (error) {
			console.error("Failed to create comment:", error);
		}
	};
	const handleCommentLike = async (commentId: number) => {
		if (!isAuthenticated) {
			navigate("/login");
			return;
		}
		try {
			await likeComment(commentId);
			setComments((prevComments) =>
				prevComments.map((comment) => {
					if (comment.id == commentId) {
						const is_liked = comment.liked_by_user;
						return {
							...comment,
							liked_by_user: !is_liked,
							likes: is_liked
								? comment.likes - 1
								: comment.likes + 1,
						};
					}
					return comment;
				})
			);
		} catch (error) {
			console.error("Failed to like comment:", error);
		}
	};
	const handlePostLike = async () => {
		if (!isAuthenticated || !postId) {
			navigate("/login");
			return;
		}
		try {
			const postIdNum = parseInt(postId, 10);

			await likePost(postIdNum);
			setPost((prevPost) => {
				if (!prevPost) return null;
				const isLiked = prevPost.liked_by_user;
				return {
					...prevPost,
					liked_by_user: !isLiked,
					likes: isLiked ? prevPost.likes - 1 : prevPost.likes + 1,
				};
			});
		} catch (error) {
			console.error("failed to like post", error);
		}
	};
	const handleDeletePost = async () => {
		if (!isAuthenticated || !postId) {
			navigate("/login");
			return;
		}
		try {
			const postIDNum = parseInt(postId, 10);
			await deletePost(postIDNum);
			navigate("/");
		} catch (error) {
			console.error("failed to delete post", error);
		}
	};
	const handleDeleteComment = async (commentID: number) => {
		if (!isAuthenticated) {
			navigate("/login");
			return;
		}
		try {
			await deleteComment(commentID);
			setComments((prevComments) =>
				prevComments.map((c) => {
					if (c.id === commentID) {
						return {
							...c,
							content: "[deleted]",
							created_by_username: "[redacted]",
						};
					}
					return c;
				})
			);
		} catch (error) {
			console.error("Failed to delete comment", error);
		}
	};
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
				<Box
					sx={{
						display: "flex",
						justifyContent: "space-between",
						alignItems: "center",
						marginTop: 2,
					}}
				>
					<Typography variant="caption">
						{post.likes || 0} Likes
					</Typography>
					<Box sx={{ display: "flex", gap: 1 }}>
						{isAuthenticated && (
							<IconButton onClick={() => handlePostLike()}>
								{post.liked_by_user ? (
									<FavoriteIcon
										sx={{
											cursor: "pointer",
											":hover": { color: "red" },
										}}
									/>
								) : (
									<FavoriteBorderIcon
										sx={{
											cursor: "pointer",
											":hover": { color: "red" },
										}}
									/>
								)}
							</IconButton>
						)}
						{isAuthenticated && user?.id == post.user_id && (
							<IconButton onClick={() => handleDeletePost()}>
								<DeleteIcon
									sx={{
										cursor: "pointer",
										":hover": { color: "red" },
									}}
								/>
							</IconButton>
						)}
					</Box>
				</Box>
			</Box>

			<Box mt={4}>
				{isAuthenticated && (
					<TextBox
						label="Add a comment..."
						onSubmit={(text) => handleSubmit(text, null)}
					/>
				)}
				{!comments || comments.length === 0 ? (
					<Typography variant="body2" color="text.secondary" mt={2}>
						No comments yet. Be the first to comment!
					</Typography>
				) : (
					comments.map((comment) => (
						<Grid key={comment.id} size={{ xs: 12, sm: 6 }}>
							<Card
								sx={{
									border: "1px solid #ddd",
									borderRadius: "8px",
									padding: "16px",
									marginBottom: "16px",
									position: "relative",
								}}
								//Todo: subreplies etc
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
								<Box>
									<Typography
										variant="caption"
										sx={{ marginTop: 2 }}
									>
										{comment.likes || 0} Likes
									</Typography>
								</Box>
								{isAuthenticated &&
									!comment.deleted &&
									comment.user_id == user?.id && (
										<IconButton
											onClick={() =>
												handleDeleteComment(comment.id)
											}
										>
											<DeleteIcon
												sx={{
													cursor: "pointer",
													":hover": { color: "red" },
												}}
											/>
										</IconButton>
									)}
								{isAuthenticated && !comment.deleted && (
									<IconButton
										sx={{
											position: "absolute",
											top: 8,
											right: 8,
										}}
										onClick={() =>
											handleCommentLike(comment.id)
										}
									>
										{comment.liked_by_user ? (
											<FavoriteIcon
												sx={{
													cursor: "pointer",
													":hover": { color: "red" },
												}}
											/>
										) : (
											<FavoriteBorderIcon
												sx={{
													cursor: "pointer",
													":hover": { color: "red" },
												}}
											/>
										)}
									</IconButton>
								)}
								<Box my={2} />
							</Card>
						</Grid>
					))
				)}
			</Box>
		</Container>
	);
};

export default PostPage;
