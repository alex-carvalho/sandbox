
val PekkoVersion = "1.0.1"

plugins {
    application
}

repositories {
    mavenCentral()
}

dependencies {
    testImplementation(libs.junit.jupiter)

    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
    implementation(platform("org.apache.pekko:pekko-bom_2.13:1.0.2"))
    implementation("ch.qos.logback:logback-classic:1.4.14")


    implementation("org.apache.pekko:pekko-actor-typed_2.13")
    implementation(libs.guava)
}

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(21))
    }
}

application {
    mainClass.set("com.ac.App")
}

tasks.named<Test>("test") {
    useJUnitPlatform()
}
