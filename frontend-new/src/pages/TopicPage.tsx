import { useEffect, useState } from "react";
import { useNavigate, useParams, Link } from "react-router-dom";
import { fetchTopicById, fetchPostsByTopicId } from "../api/forum";
import type { Topic, Post } from "../types/models";
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
	const [topic, setTopic] = useState<Topic | null>(null);
	const [loading, setLoading] = useState(true);
	const [posts, setPosts] = useState<Post[]>([]);
	const [error, setError] = useState<string | null>(null);
	const { isAuthenticated, user } = useAuth();

	useEffect(() => {
		const loadTopicAndPosts = async () => {
			if (!topicId) return;
			try {
				setLoading(true);
				const topicIdNum = parseInt(topicId, 10);
				if (isNaN(topicIdNum)) {
					setError("Invalid Topic ID");
					setLoading(false);
					return;
				}
				const topicData = await fetchTopicById(topicIdNum);
				const postsData = await fetchPostsByTopicId(topicIdNum);
				setPosts(postsData);
				setTopic(topicData);
			} catch (error) {
				console.error("Failed to fetch topic or posts:", error);
				setError("Failed to load topic. It may not exist.");
			} finally {
				setLoading(false);
			}
		};

		loadTopicAndPosts();
	}, [topicId]);
	if (loading) {
		return (
			<Box display="flex" justifyContent="center" mt={4}>
				<CircularProgress />
			</Box>
		);
	}
	if (error || !topic) {
		return (
			<Container sx={{ mt: 4 }}>
				<Alert severity="error">{error || "Topic not found"}</Alert>
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
					{topic.title}
				</Typography>
				<Typography
					variant="subtitle1"
					color="text.secondary"
					gutterBottom
				>
					{topic.description}
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
					<Typography variant="h6">Posts</Typography>
					{isAuthenticated && (
						<Button
							variant="contained"
							color="primary"
							onClick={() =>
								navigate(`/topics/${topicId}/create`)
							}
						>
							Create Post
						</Button>
					)}
				</Box>
				{!posts || posts.length === 0 ? (
					<Typography variant="body2" color="text.secondary" mt={2}>
						No posts yet. Be the first to post!
					</Typography>
				) : (
					posts.map((post) => (
						<Grid key={post.id} size={{ xs: 12, sm: 6 }}>
							<CardActionArea
								sx={{
									border: "1px solid #ddd",
									borderRadius: "8px",
									padding: "16px",
									marginBottom: "16px",
								}}
								component={Link}
								to={`/topics/${topicId}/posts/${post.id}`}
							>
								<Typography variant="h6">
									{post.title}
								</Typography>
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
									{post.content}
								</Typography>
								<Typography
									variant="caption"
									color="text.secondary"
								>
									Posted by • {user?.username} •{" "}
									{timeAgo(post.created_at)} (on{" "}
									{formatDate(post.created_at)})
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
