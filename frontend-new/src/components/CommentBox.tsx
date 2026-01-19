import { Card, Typography, Box, IconButton, Grid, Button } from "@mui/material";
import FavoriteIcon from "@mui/icons-material/Favorite";
import FavoriteBorderIcon from "@mui/icons-material/FavoriteBorder";
import DeleteIcon from "@mui/icons-material/Delete";
import { timeAgo, formatDate } from "../utils/date";
import type { CommentNode } from "../utils/commentTree";
import { useState } from "react";
import TextBox from "./TextBox";
import EditIcon from "@mui/icons-material/Edit";
import AddCommentIcon from "@mui/icons-material/AddComment";
interface CommentProps {
	comment: CommentNode;
	currentUserId?: number;
	isAuthenticated: boolean;
	onLike: (id: number) => void;
	onDelete: (id: number) => void;
	onEdit: (id: number, content: string) => void;
	onReply: (content: string, parentID: number) => void;
}

const CommentBox = ({ comment, onReply, ...props }: CommentProps) => {
	const [isDeleted, setIsDeleted] = useState(comment.deleted);
	const isUpdated = comment.updated_at !== comment.created_at;
	const [isEditing, setIsEditing] = useState(false);
	const [content, setContent] = useState(comment.content);
	const [isReplying, setIsReplying] = useState(false);

	const handleReplySubmit = (content: string) => {
		onReply(content, comment.id);
		setIsReplying(false);
	};
	return (
		<Box sx={{ marginBottom: 2 }}>
			<Grid size={{ xs: 12, sm: 6 }}>
				<Card
					sx={{
						border: "1px solid #ddd",
						borderRadius: "8px",
						padding: "16px",
						marginBottom: "16px",
						position: "relative",
						backgroundColor: isDeleted
							? "#f5f5f5"
							: "var(--color-platinum-100)",
					}}
				>
					{isEditing ? (
						<>
							<TextBox
								onSubmit={(content) => {
									props.onEdit(comment.id, content);
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
								{!isDeleted ? comment.content : "[deleted]"}
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
											: comment.created_at,
									)}
									)
								</Typography>
							)}

							<Box mt={2}>
								{!isDeleted && (
									<Typography variant="caption">
										{comment.likes || 0} Likes
									</Typography>
								)}
							</Box>

							{props.isAuthenticated && !isDeleted && (
								<Box
									sx={{
										position: "absolute",
										bottom: 8,
										right: 8,
										display: "flex",
										gap: 1,
									}}
								>
									<IconButton
										onClick={() => props.onLike(comment.id)}
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

									{props.currentUserId ===
										comment.user_id && (
										<>
											<IconButton
												onClick={() =>
													setIsEditing(true)
												}
											>
												<EditIcon
													sx={{
														cursor: "pointer",
														":hover": {
															color: "red",
														},
													}}
												/>
											</IconButton>
											<IconButton
												onClick={() => {
													props.onDelete(comment.id);
													setIsDeleted(true);
												}}
											>
												<DeleteIcon
													sx={{
														cursor: "pointer",
														":hover": {
															color: "red",
														},
													}}
												/>
											</IconButton>
											<IconButton
												onClick={() =>
													setIsReplying(!isReplying)
												}
											>
												<AddCommentIcon />
											</IconButton>
										</>
									)}
								</Box>
							)}
						</>
					)}
				</Card>
				{isReplying && (
					<Box sx={{ ml: 4, mt: 2 }}>
						<TextBox
							onSubmit={handleReplySubmit}
							label="Write a reply..."
						/>
					</Box>
				)}
			</Grid>
			{comment.children && comment.children.length > 0 && (
				<Box sx={{ ml: 4, mt: 2, borderLeft: "2px solid #eee", pl: 2 }}>
					{comment.children.map((child) => (
						<CommentBox
							key={child.id}
							comment={child}
							onReply={onReply}
							{...props}
						/>
					))}
				</Box>
			)}
		</Box>
	);
};

export default CommentBox;
