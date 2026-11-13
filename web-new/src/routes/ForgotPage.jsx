import ForgotForm from "../schemas/ForgotForm";
import { createResource, createSignal } from "solid-js";
//import { challengeService } from "../services";

function ForgotPage() {
		const [sent, setSent] = createSignal(false);

		const sendEmail = (data) => {
				const email = data.email;
				setSent(true);
		}
		
		return (
				<>
						<Show when={sent()}
									fallback={
											<>
													<h3>SÄkhöpoisti!!</h3>
													<ForgotForm
															onSubmit={sendEmail}
													/>
											</>
									}>
								Sähköposti lähetetty, moi
						</Show>
				</>
		);
}

export default ForgotPage;