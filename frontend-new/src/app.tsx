import { BrowserRouter, Routes, Route } from "react-router-dom";
import ResponsiveAppBar from "./components/NavBar";
import LoginPage from "./pages/LoginPage";
import HomePage from "./pages/HomePage";
import TopicPage from "./pages/TopicPage";
import PostPage from "./pages/PostPage";
import CreatePostPage from "./pages/CreatePostPage";
import CreateTopicsPage from "./pages/CreateTopicsPage";
import ProtectedRoutes from "./components/ProtectedRoutes";
import RegisterPage from "./pages/RegisterPage";

function App() {
	return (
		<BrowserRouter>
			<ResponsiveAppBar />
			<Routes>
				<Route path="/" element={<HomePage />} />
				<Route path="/login" element={<LoginPage />} />
				<Route path="/register" element={<RegisterPage />} />
				<Route path="/topics/:topicId" element={<TopicPage />} />
				<Route
					path="/topics/:topicId/posts/:postId"
					element={<PostPage />}
				/>
				<Route element={<ProtectedRoutes />}>
					<Route
						path="/topics/:topicId/create"
						element={<CreatePostPage />}
					/>
					<Route
						path="/topics/create"
						element={<CreateTopicsPage />}
					/>
				</Route>
				<Route path="*" element={<HomePage />} />
			</Routes>
		</BrowserRouter>
	);
}
export default App;
