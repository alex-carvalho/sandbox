package com.ac.batch.csvtomongo.job;

import org.springframework.util.StringUtils;

import java.beans.PropertyEditorSupport;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;

public class LocalDatePropertyEditor extends PropertyEditorSupport {

    private final String pattern;

    LocalDatePropertyEditor(String pattern) {
        this.pattern = pattern;
    }

    @Override
    public String getAsText() {
        LocalDate value = (LocalDate) getValue();
        return (value != null ? value.toString() : "");
    }

    @Override
    public void setAsText(String text) throws IllegalArgumentException {
        if (StringUtils.isEmpty(text)) return;

        setValue(LocalDate.parse(text, DateTimeFormatter.ofPattern(pattern)));
    }
}
