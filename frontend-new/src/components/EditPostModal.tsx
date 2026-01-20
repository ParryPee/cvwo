import { useState, useEffect } from "react";
import {
	Dialog,
	DialogTitle,
	DialogContent,
	DialogActions,
	TextField,
	Button,
	Box,
} from "@mui/material";

interface EditPostModalProps {
	open: boolean;
	onClose: () => void;
	onSubmit: (newTitle: string, newContent: string) => Promise<void>;
	initialTitle: string;
	initialContent: string;
}

const EditPostModal = ({
	open,
	onClose,
	onSubmit,
	initialTitle,
	initialContent,
}: EditPostModalProps) => {
	const [title, setTitle] = useState(initialTitle);
	const [content, setContent] = useState(initialContent);
	const [loading, setLoading] = useState(false);

	useEffect(() => {
		if (open) {
			setContent(initialContent);
		}
	}, [open, initialContent]);

	const handleSubmit = async () => {
		if (content.trim() === "") {
			return;
		}
		try {
			setLoading(true);
			await onSubmit(title, content);
		} catch (error) {
			console.error("Failed to update post", error);
		} finally {
			setLoading(false);
		}
	};
	return (
		<Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
			<Box
				sx={{
					bgcolor: "var(--color-lavender-grey-300)",
				}}
			>
				<DialogTitle>Edit Post</DialogTitle>
				<DialogContent>
					<TextField
						autoFocus
						margin="dense"
						label="Content"
						fullWidth
						multiline
						rows={6}
						variant="outlined"
						value={title}
						onChange={(e) => setTitle(e.target.value)}
						disabled={loading}
					/>
					<TextField
						autoFocus
						margin="dense"
						label="Content"
						fullWidth
						multiline
						rows={6}
						variant="outlined"
						value={content}
						onChange={(e) => setContent(e.target.value)}
						disabled={loading}
					/>
				</DialogContent>
				<DialogActions>
					<Button
						onClick={onClose}
						disabled={loading}
						sx={{
							bgcolor: "var(--color-lavender-grey-500)",
							color: "white",
						}}
					>
						Cancel
					</Button>
					<Button
						onClick={handleSubmit}
						variant="contained"
						color="primary"
						disabled={loading || !content.trim()}
						sx={{
							bgcolor: "var(--color-lavender-grey-500)",
						}}
					>
						{loading ? "Saving..." : "Save"}
					</Button>
				</DialogActions>
			</Box>
		</Dialog>
	);
};

export default EditPostModal;
