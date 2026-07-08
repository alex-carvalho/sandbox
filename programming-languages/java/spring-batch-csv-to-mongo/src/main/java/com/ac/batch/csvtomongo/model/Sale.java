package com.ac.batch.csvtomongo.model;

import org.springframework.data.mongodb.core.mapping.Document;

import java.time.LocalDate;

@Document(collection = "sales")
public class Sale {

    private String region;
    private String country;
    private String itemType;
    private String salesChannel;
    private String orderPriority;
    private LocalDate orderDate;
    private String orderID;
    private LocalDate shipDate;
    private Long unitsSold;
    private Double unitPrice;
    private Double unitCost;
    private Double totalRevenue;
    private Double totalCost;
    private Double totalProfit;


    public String getRegion() {
        return region;
    }

    public void setRegion(String region) {
        this.region = region;
    }

    public String getCountry() {
        return country;
    }

    public void setCountry(String country) {
        this.country = country;
    }

    public String getItemType() {
        return itemType;
    }

    public void setItemType(String itemType) {
        this.itemType = itemType;
    }

    public String getSalesChannel() {
        return salesChannel;
    }

    public void setSalesChannel(String salesChannel) {
        this.salesChannel = salesChannel;
    }

    public String getOrderPriority() {
        return orderPriority;
    }

    public void setOrderPriority(String orderPriority) {
        this.orderPriority = orderPriority;
    }

    public LocalDate getOrderDate() {
        return orderDate;
    }

    public void setOrderDate(LocalDate orderDate) {
        this.orderDate = orderDate;
    }

    public String getOrderID() {
        return orderID;
    }

    public void setOrderID(String orderID) {
        this.orderID = orderID;
    }

    public LocalDate getShipDate() {
        return shipDate;
    }

    public void setShipDate(LocalDate shipDate) {
        this.shipDate = shipDate;
    }

    public Long getUnitsSold() {
        return unitsSold;
    }

    public void setUnitsSold(Long unitsSold) {
        this.unitsSold = unitsSold;
    }

    public Double getUnitPrice() {
        return unitPrice;
    }

    public void setUnitPrice(Double unitPrice) {
        this.unitPrice = unitPrice;
    }

    public Double getUnitCost() {
        return unitCost;
    }

    public void setUnitCost(Double unitCost) {
        this.unitCost = unitCost;
    }

    public Double getTotalRevenue() {
        return totalRevenue;
    }

    public void setTotalRevenue(Double totalRevenue) {
        this.totalRevenue = totalRevenue;
    }

    public Double getTotalCost() {
        return totalCost;
    }

    public void setTotalCost(Double totalCost) {
        this.totalCost = totalCost;
    }

    public Double getTotalProfit() {
        return totalProfit;
    }

    public void setTotalProfit(Double totalProfit) {
        this.totalProfit = totalProfit;
    }

    @Override
    public String toString() {
        return "Sale{" +
                "region='" + region + '\'' +
                ", country='" + country + '\'' +
                ", itemType='" + itemType + '\'' +
                ", salesChannel='" + salesChannel + '\'' +
                ", orderPriority='" + orderPriority + '\'' +
                ", orderDate='" + orderDate + '\'' +
                ", orderID='" + orderID + '\'' +
                ", shipDate='" + shipDate + '\'' +
                ", unitsSold='" + unitsSold + '\'' +
                ", unitPrice='" + unitPrice + '\'' +
                ", unitCost='" + unitCost + '\'' +
                ", totalRevenue='" + totalRevenue + '\'' +
                ", totalCost='" + totalCost + '\'' +
                ", totalProfit='" + totalProfit + '\'' +
                '}';
    }
}
