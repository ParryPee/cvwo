import * as React from "react";
import { useNavigate, Link, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import Menu from "@mui/material/Menu";
import MenuIcon from "@mui/icons-material/Menu";
import Container from "@mui/material/Container";
import Avatar from "@mui/material/Avatar";
import Button from "@mui/material/Button";
import Tooltip from "@mui/material/Tooltip";
import MenuItem from "@mui/material/MenuItem";
import Logo from "../../logo.png";

const pages = [
	{ name: "Home", link: "/" },
	{ name: "Explore", link: "/explore" },
];

function ResponsiveAppBar() {
	const { user, isAuthenticated, logout } = useAuth();
	const navigate = useNavigate();
	const location = useLocation();

	const [anchorElNav, setAnchorElNav] = React.useState<null | HTMLElement>(
		null,
	);
	const [anchorElUser, setAnchorElUser] = React.useState<null | HTMLElement>(
		null,
	);
	const [searchTerm, setSearchTerm] = React.useState("");

	const handleSearch = (e: React.KeyboardEvent) => {
		if (e.key === "Enter") {
			navigate(`/?q=${searchTerm}`);
		}
	};

	const handleOpenNavMenu = (event: React.MouseEvent<HTMLElement>) => {
		setAnchorElNav(event.currentTarget);
	};
	const handleOpenUserMenu = (event: React.MouseEvent<HTMLElement>) => {
		setAnchorElUser(event.currentTarget);
	};

	const handleCloseNavMenu = () => {
		setAnchorElNav(null);
	};

	const handleCloseUserMenu = () => {
		setAnchorElUser(null);
	};

	const handleLogout = () => {
		handleCloseUserMenu();
		logout();
		navigate("/");
	};

	return (
		<AppBar position="static">
			<Container maxWidth={false} className="bg-lavender-grey-700">
				<Toolbar disableGutters>
					<img
						src={Logo}
						alt="CVWO Logo"
						style={{
							width: "auto",
							marginRight: "10px",
							maxWidth: "100px",
						}}
					/>
					<Typography
						variant="h6"
						noWrap
						component={Link}
						to="/"
						sx={{
							mr: 2,
							display: { xs: "none", md: "flex" },
							fontFamily: "monospace",
							fontWeight: 700,
							letterSpacing: ".3rem",
							color: "inherit",
							textDecoration: "none",
						}}
					>
						CVWO
					</Typography>

					<Box
						sx={{
							flexGrow: 1,
							display: { xs: "flex", md: "none" },
						}}
					>
						<IconButton
							size="large"
							aria-label="account of current user"
							aria-controls="menu-appbar"
							aria-haspopup="true"
							onClick={handleOpenNavMenu}
							color="inherit"
						>
							<MenuIcon />
						</IconButton>
						<Menu
							id="menu-appbar"
							anchorEl={anchorElNav}
							anchorOrigin={{
								vertical: "bottom",
								horizontal: "left",
							}}
							keepMounted
							transformOrigin={{
								vertical: "top",
								horizontal: "left",
							}}
							open={Boolean(anchorElNav)}
							onClose={handleCloseNavMenu}
							sx={{ display: { xs: "block", md: "none" } }}
						>
							{pages.map((page) => (
								<MenuItem
									key={page.name}
									component={Link}
									to={page.link}
									onClick={handleCloseNavMenu}
								>
									<Typography textAlign="center">
										{page.name}
									</Typography>
								</MenuItem>
							))}
						</Menu>
					</Box>
					<Typography
						variant="h5"
						noWrap
						component="a"
						href="/"
						sx={{
							mr: 2,
							display: { xs: "flex", md: "none" },
							flexGrow: 1,
							fontFamily: "monospace",
							fontWeight: 700,
							letterSpacing: ".3rem",
							color: "inherit",
							textDecoration: "none",
						}}
					>
						CVWO
					</Typography>

					<Box
						sx={{
							flexGrow: 1,
							display: { xs: "none", md: "flex" },
						}}
					>
						{pages.map((page) => (
							<Button
								key={page.name}
								component={Link}
								to={page.link}
								onClick={handleCloseNavMenu}
								sx={{
									my: 2,
									color: "white",
									display: "block",
									mr: 2,
									borderRadius: 2,
									"&:hover": {
										bgcolor: "var(--color-flag-red-500)",
									},
								}}
							>
								{page.name}
							</Button>
						))}
					</Box>
					{location.pathname !== "/" && (
						<Box
							sx={{
								mr: 2,
								display: "flex",
								alignItems: "center",
								backgroundColor: "var(--color-platinum-100)",
								borderRadius: 4,
								paddingX: 1,
							}}
						>
							<input
								placeholder="Search..."
								value={searchTerm}
								onChange={(e) => setSearchTerm(e.target.value)}
								onKeyDown={handleSearch}
								style={{
									padding: "8px",
									borderRadius: "4px",
									border: "none",
									outline: "none",
									color: "var(--color-lavender-grey-900)",
								}}
							/>
						</Box>
					)}

					<Box sx={{ flexGrow: 0 }}>
						{!isAuthenticated ? (
							<Button
								color="inherit"
								component={Link}
								to="/login"
								sx={{
									borderColor: "white",
									color: "white",
									"&:hover": {
										bgcolor: "var(--color-flag-red-600)",
									},
									borderRadius: 4,
								}}
							>
								Login
							</Button>
						) : (
							<>
								<Tooltip title="Open settings">
									<IconButton
										onClick={handleOpenUserMenu}
										sx={{ p: 0 }}
									>
										<Avatar alt={user?.username || "User"}>
											{user?.username
												? user.username
														.charAt(0)
														.toUpperCase()
												: "U"}
										</Avatar>
									</IconButton>
								</Tooltip>
								<Menu
									sx={{ mt: "45px" }}
									id="menu-appbar"
									anchorEl={anchorElUser}
									anchorOrigin={{
										vertical: "top",
										horizontal: "right",
									}}
									keepMounted
									transformOrigin={{
										vertical: "top",
										horizontal: "right",
									}}
									open={Boolean(anchorElUser)}
									onClose={handleCloseUserMenu}
								>
									<MenuItem onClick={handleLogout}>
										<Typography
											textAlign="center"
											color="error"
										>
											Logout
										</Typography>
									</MenuItem>
								</Menu>
							</>
						)}
					</Box>
				</Toolbar>
			</Container>
		</AppBar>
	);
}
export default ResponsiveAppBar;
