import { createSignal } from "solid-js";
import trashCan from "../../assets/can.svg";

function ConfirmButton(props) {
		const [confirm, setConfirm] = createSignal(false);

		return (
				<Show when={confirm()}
							fallback={
									<button class="btn-icon"
													onClick={() => setConfirm(true)}>
											<img src={trashCan} alt="Poista" />
									</button>
							}>
						<button class="btn-danger" onClick={props.onClick}>
								Varmista poisto
						</button>
						<button class="btn-secondary"
										onClick={() => setConfirm(false)}>
								Peru
						</button>
				</Show>
		)
}

export default ConfirmButton;