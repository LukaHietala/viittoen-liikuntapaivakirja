import { createForm, Field, Form, reset, setInput, getInput } from "@formisch/solid";
import * as v from "valibot";

const performanceSchema = v.pipe(
		v.object({
				points: v.pipe(
						v.toNumber(),
						v.number(),
						v.integer(),
						v.toMinValue(1)
				),
		})
);

function PerformanceForm(props) {
		const cancel = () => props.cancel;

		const handleSubmit = async (data) => {
				await props.onSubmit(data);
		};

		const performanceForm = createForm({
				schema: performanceSchema,
		});

		return (
				<Form
						of={performanceForm}
						onSubmit={handleSubmit}
						class="form"
				>				
						<Field
								of={performanceForm}
								path={["points"]}
						>
								{(field) => (
										<div>
												<label for={field.props.name}>Pisteitä:</label>
												<br />
												<input
														{...field.props}
														id={field.props.name}
														value={field.input}
														type="number"
														required
												/>
												{field.errors && <small class="form-error">{field.errors[0]}</small>}
										</div>
								)}
						</Field>
						
						<button
								class="btn-secondary"
								type="button"
								onClick={() => cancel()}
						>
								Peruuta
						</button>
						<button
								class="btn-primary"
								type="submit"
						>
								Tallenna
						</button>
				</Form>
		);
}

export default PerformanceForm;