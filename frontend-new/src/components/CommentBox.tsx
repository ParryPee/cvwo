import {
	Card,
	Typography,
	Box,
	IconButton,
	Grid,
	TextField,
	Button,
} from "@mui/material";
import FavoriteIcon from "@mui/icons-material/Favorite";
import FavoriteBorderIcon from "@mui/icons-material/FavoriteBorder";
import DeleteIcon from "@mui/icons-material/Delete";
import { timeAgo, formatDate } from "../utils/date";
import type { Comment } from "../types/models";
import { useState } from "react";
import TextBox from "./TextBox";
import EditIcon from "@mui/icons-material/Edit";

interface CommentProps {
	comment: Comment;
	currentUserId?: number;
	isAuthenticated: boolean;
	onLike: (id: number) => void;
	onDelete: (id: number) => void;
	onEdit: (id: number, content: string) => void;
}

const CommentBox = ({
	comment,
	currentUserId,
	isAuthenticated,
	onLike,
	onDelete,
	onEdit,
}: CommentProps) => {
	const isDeleted = comment.deleted;
	const isUpdated = comment.updated_at !== comment.created_at;
	console.log(isUpdated);
	const [isEditing, setIsEditing] = useState(false);
	const [content, setContent] = useState(comment.content);

	return (
		<Grid size={{ xs: 12, sm: 6 }}>
			<Card
				sx={{
					border: "1px solid #ddd",
					borderRadius: "8px",
					padding: "16px",
					marginBottom: "16px",
					position: "relative",
					backgroundColor: isDeleted ? "#f5f5f5" : "white",
				}}
			>
				{isEditing ? (
					<>
						<TextBox
							onSubmit={(content) => {
								onEdit(comment.id, content);
								setContent(content);
								setIsEditing(false);
							}}
							label="content"
							content={content}
						/>
						<Button
							variant="text"
							sx={{ color: "red" }}
							onClick={() => setIsEditing(false)}
						>
							Cancel
						</Button>
					</>
				) : (
					<>
						<Typography
							variant="body1"
							sx={{
								fontStyle: isDeleted ? "italic" : "normal",
								color: isDeleted
									? "text.secondary"
									: "text.primary",
								display: "-webkit-box",
								WebkitLineClamp: 2,
								WebkitBoxOrient: "vertical",
								overflow: "hidden",
								textOverflow: "ellipsis",
							}}
						>
							{content}
						</Typography>

						{!isDeleted && (
							<Typography
								variant="caption"
								color="text.secondary"
							>
								{isUpdated ? "Updated" : "Posted"} by •{" "}
								{comment.created_by_username} •{" "}
								{isUpdated
									? timeAgo(comment.updated_at)
									: timeAgo(comment.created_at)}{" "}
								(on{" "}
								{formatDate(
									isUpdated
										? comment.updated_at
										: comment.created_at
								)}
								)
							</Typography>
						)}

						<Box mt={2}>
							<Typography variant="caption">
								{comment.likes || 0} Likes
							</Typography>
						</Box>

						{isAuthenticated && !isDeleted && (
							<>
								<IconButton
									sx={{
										position: "absolute",
										top: 8,
										right: 8,
									}}
									onClick={() => onLike(comment.id)}
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

								{currentUserId === comment.user_id && (
									<>
										<IconButton
											onClick={() => setIsEditing(true)}
										>
											<EditIcon
												sx={{
													cursor: "pointer",
													":hover": { color: "red" },
												}}
											/>
										</IconButton>
										<IconButton
											sx={{
												position: "absolute",
												bottom: 8,
												right: 8,
											}}
											onClick={() => onDelete(comment.id)}
										>
											<DeleteIcon
												sx={{
													cursor: "pointer",
													":hover": { color: "red" },
												}}
											/>
										</IconButton>
									</>
								)}
							</>
						)}
					</>
				)}
			</Card>
		</Grid>
	);
};

export default CommentBox;
