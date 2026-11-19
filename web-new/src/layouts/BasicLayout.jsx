import { useMatch } from "@solidjs/router";

function BasicLayout(props) {
		return (
				<>
						<main class="container">
								{props.children}
						</main>
				</>
		);
}

export default BasicLayout;