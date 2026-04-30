package com.ac.pulsar;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.pulsar.core.PulsarAdministration;
import org.springframework.pulsar.core.PulsarTemplate;
import org.springframework.pulsar.core.PulsarTopic;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.verify;

@ExtendWith(MockitoExtension.class)
class PulsarApplicationTests {

	@Mock
	private PulsarTemplate<String> pulsarTemplate;

	@Mock
	private PulsarAdministration pulsarAdministration;

	@Test
	void forwardsOrdersToPaymentTopic() {
		var consumer = new OrderToPaymentConsumer(this.pulsarTemplate, "payment");

		consumer.handle("order-123");

		verify(this.pulsarTemplate).send("payment", "order-123");
	}

	@Test
	void createsOrdersAndPaymentTopicsOnStartup() throws Exception {
		var configuration = new PulsarTopicConfiguration();
		var topicsCaptor = ArgumentCaptor.forClass(PulsarTopic[].class);

		configuration.pulsarTopicInitializer(this.pulsarAdministration, "orders", "payment").run(null);

		verify(this.pulsarAdministration).createOrModifyTopics(topicsCaptor.capture());
		assertThat(topicsCaptor.getValue()).hasSize(2);
		assertThat(topicsCaptor.getValue()[0].topicName()).endsWith("/orders");
		assertThat(topicsCaptor.getValue()[1].topicName()).endsWith("/payment");
	}

}
